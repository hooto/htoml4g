package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("invalid command")
	}

	files := []string{}
	if strings.HasSuffix(os.Args[1], ".go") {
		files = append(files, os.Args[1])
	} else {
		out, err := exec.Command("find", os.Args[1], "-type", "f", "-name", "*.go").Output()
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
	fmt.Println(files)

	for _, file := range files {
		if err := filter(file); err != nil {
			log.Fatal(err)
		}
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

		if n1 := bytes.IndexByte(line, '`'); n1 > 0 {
			if n2 := bytes.Index(line, []byte("toml:\"")); n2 <= n1 {
				if nj1 := bytes.Index(line, []byte("json:\"")); nj1 > n1 {
					if tl := bytes.IndexByte(line[nj1+6:], '"'); tl > 0 {
						line = []byte(string(line[:nj1+7+tl]) +
							" toml:\"" + string(line[nj1+6:nj1+6+tl]) + "\"" +
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
