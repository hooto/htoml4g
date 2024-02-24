// Copyright 2019 Eryx <evorui аt gmail dοt com>, All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

var (
	nameRes = []*regexp.Regexp{
		regexp.MustCompile(`,name=(.*?),`),
		regexp.MustCompile(`protobuf_oneof:"(.*?)"`),
	}
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("invalid command")
	}

	files := []string{}
	for _, v := range strings.Split(os.Args[1], ",") {
		if strings.HasSuffix(v, ".go") {
			files = append(files, v)
		} else {
			out, err := exec.Command("find", v, "-type", "f", "-name", "*.go").Output()
			if err != nil {
				log.Fatal(err)
			}
			arr := strings.Split(string(out), "\n")
			for _, v := range arr {
				if strings.HasSuffix(v, ".go") {
					files = append(files, v)
				}
			}
		}
	}

	for _, file := range files {
		if err := filter(file); err != nil {
			log.Fatal(err)
		}
		fmt.Println("  ok", file)
	}
}

func filter(f string) error {

	fp, err := os.Open(f)
	if err != nil {
		return err
	}

	var (
		r   = bufio.NewReader(fp)
		chg = false
		buf bytes.Buffer
	)

	for {

		line, _, err := r.ReadLine()
		if err != nil {
			break
		}

		nameHit := func(v string) string {
			for _, nameRe := range nameRes {
				if mat := nameRe.FindAllStringSubmatch(v, 1); len(mat) > 0 {
					return mat[0][1]
				}
			}
			return ""
		}

		if n1 := bytes.IndexByte(line, '`'); n1 > 0 {
			if n2 := bytes.Index(line, []byte("toml:\"")); n2 <= n1 {

				if fieldName := nameHit(string(line[n1:])); fieldName != "" {

					for _, s := range []string{
						"protobuf_oneof:\"",
						"proto3,oneof",
					} {
						if nj1 := bytes.Index(line, []byte(s)); nj1 > n1 {
							if tl := bytes.IndexByte(line[nj1+len(s):], '"'); tl >= 0 {
								line = []byte(string(line[:nj1+len(s)+1+tl]) +
									" json:\"" + fieldName + "\"" +
									string(line[nj1+len(s)+1+tl:]))
								chg = true
							}
						}
					}
				}

				if nj1 := bytes.Index(line, []byte("json:\"")); nj1 > n1 {
					if tl := bytes.IndexByte(line[nj1+6:], '"'); tl > 0 {
						line = []byte(string(line[:nj1+7+tl]) +
							" toml:\"" + string(line[nj1+6:nj1+6+tl]) + "\"" +
							" yaml:\"" + string(line[nj1+6:nj1+6+tl]) + "\"" +
							string(line[nj1+7+tl:]))
						chg = true
					}
				}
			}
		}

		buf.Write(line)
		buf.Write([]byte("\n"))
	}

	fp.Close()

	if chg {
		fpo, err := os.OpenFile(f, os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		fpo.Seek(0, 0)
		fpo.Truncate(0)
		fpo.Write(buf.Bytes())
		fpo.Sync()

		fpo.Close()
	}

	return nil
}
