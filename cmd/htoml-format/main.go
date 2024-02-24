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
	"regexp"
	"strings"

	"github.com/hooto/htoml4g/internal/toml"
)

func main() {

	if len(os.Args) < 2 {
		return
	}

	if !strings.HasSuffix(strings.ToLower(os.Args[1]), ".toml") {
		return
	}

	filter(os.Args[1])
}

func filter(f string) error {

	fp, err := os.Open(f)
	if err != nil {
		return err
	}
	b, err := io.ReadAll(fp)
	fp.Close()
	if err != nil {
		return err
	}

	var w bytes.Bufer

	fpo, err := os.OpenFile(f, os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	fpo.Seek(0, 0)
	fpo.Truncate(0)
	fpo.Write(buf.Bytes())
	fpo.Sync()

	fpo.Close()

	enc := toml.NewEncoder(w)
	return enc.Encode(v)

	return nil
}
