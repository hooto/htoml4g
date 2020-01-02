// Copyright 2019 Eryx <evorui аt gmаil dοt cοm>, All rights reserved.
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

package hconf

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	Version = "0.9.0"
)

func Decode(obj interface{}, bs []byte) error {

	if _, err := toml.Decode(string(bs), obj); err != nil {
		return err
	}

	return nil
}

func DecodeFromFile(obj interface{}, file string) error {

	fp, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fp.Close()

	bs, err := ioutil.ReadAll(fp)
	if err != nil {
		return err
	}

	if _, err := toml.Decode(string(bs), obj); err != nil {
		return err
	}

	return nil
}

func Encode(obj interface{}) ([]byte, error) {

	var (
		buf bytes.Buffer
		enc = toml.NewEncoder(&buf)
	)

	if err := enc.Encode(obj); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func EncodeToFile(obj interface{}, file string) error {

	var (
		buf bytes.Buffer
		enc = toml.NewEncoder(&buf)
	)

	if err := enc.Encode(obj); err != nil {
		return err
	}

	fpo, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer fpo.Close()

	fpo.Seek(0, 0)
	fpo.Truncate(0)

	wbuf := bufio.NewWriter(fpo)

	for {

		line, err := buf.ReadBytes('\n')
		if err != nil {
			break
		}

		if len(line) > 80 && line[len(line)-2] == '"' && line[len(line)-3] != '"' {

			if n := bytes.Index(line, []byte(" = \"")); n > 0 {
				if nb := bytes.Index(line[n+4:len(line)-2], []byte("\\n")); nb >= 0 {
					wbuf.Write(line[:n+4])
					wbuf.Write([]byte(`""`))
					wbuf.Write(bytes.Replace(line[n+4:len(line)-2], []byte("\\n"), []byte("\n"), -1))
					wbuf.Write([]byte("\"\"\"\n"))
					continue
				}
			}
		}

		wbuf.Write(line)
	}

	return wbuf.Flush()
}
