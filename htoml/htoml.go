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

package htoml

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"

	"github.com/hooto/htoml4g/internal/toml" // fork from "github.com/BurntSushi/toml"
)

const (
	Version = "0.9.4"
)

// EncodeOptions sets options such as indent.
type EncodeOptions struct {
	Indent string
}

// NewEncodeOptions returns a new EncodeOptions with default values.
func NewEncodeOptions() *EncodeOptions {
	return &EncodeOptions{Indent: ""}
}

// Decode parses the TOML-encoded data and stores the result in the value pointed to by v.
func Decode(bs []byte, v interface{}) error {

	if _, err := toml.Decode(string(bs), v); err != nil {
		return err
	}

	return nil
}

// DecodeFromFile parses the TOML-encoded data from file and stores the result in the value pointed to by v.
func DecodeFromFile(file string, v interface{}) error {

	fp, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fp.Close()

	bs, err := ioutil.ReadAll(fp)
	if err != nil {
		return err
	}

	if _, err = toml.Decode(string(bs), v); err != nil {
		return err
	}

	return nil
}

func encodeOptionsParse(args ...interface{}) *EncodeOptions {
	var opts *EncodeOptions
	for _, arg := range args {
		switch arg.(type) {
		case EncodeOptions:
			o := arg.(EncodeOptions)
			opts = &o

		case *EncodeOptions:
			opts = arg.(*EncodeOptions)
		}
	}
	if opts == nil {
		opts = NewEncodeOptions()
	}
	return opts
}

// Encode returns the TOML encoding of v.
func Encode(v interface{}, args ...interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := prettyEncode(v, &buf, encodeOptionsParse(args...)); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// EncodeToFile writes the TOML encoding to file.
func EncodeToFile(v interface{}, file string, args ...interface{}) error {

	fp, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0640)
	if err != nil {
		return err
	}
	defer fp.Close()

	fp.Seek(0, 0)
	fp.Truncate(0)

	w := bufio.NewWriter(fp)

	if err = prettyEncode(v, w, encodeOptionsParse(args...)); err != nil {
		return err
	}

	return w.Flush()
}

func prettyEncode(v interface{}, w io.Writer, opts *EncodeOptions) error {
	enc := toml.NewEncoder(w)
	if opts != nil {
		enc.Indent = opts.Indent
	}
	return enc.Encode(v)
}
