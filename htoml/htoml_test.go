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
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

type tomlConfig struct {
	Title   string
	Owner   ownerInfo
	DB      database `json:"database" toml:"database"`
	Servers []server
	Clients clients
}

type ownerInfo struct {
	Name string
	Org  string `json:"organization" toml:"organization" desc:"Org Description"`
	Bio  string
}

type database struct {
	Server  string
	Ports   []int
	ConnMax int `json:"connection_max" toml:"connection_max"`
	Enabled bool
}

type richStruct struct {
	Name   string `toml:"name"`
	Driver func(string) error
}

type server struct {
	IP string
	DC string
}

type clients struct {
	Hosts []string
}

var (
	tomlRaw = []byte(`
title = "TOML Example"

[owner]
name = "Test Robot"
organization = "Demo"
bio = "Demo Tester"

[database]
server = "192.168.1.1"
ports = [ 8001, 8001, 8002 ]
connection_max = 5000
enabled = true

[[servers]]
ip = "10.0.0.1"
dc = "eqdc10"

[[servers]]
ip = "10.0.0.2"
dc = "eqdc10"

[clients]
hosts = [
  "alpha",
  "omega"
]`)
	jsonRaw []byte
	err     error
	objRaw  tomlConfig
)

func init() {
	var item tomlConfig
	if err = Decode(tomlRaw, &item); err != nil {
		panic(err)
	}
	if jsonRaw, err = json.Marshal(item); err != nil {
		panic(err)
	}
	objRaw = item
}

func Test_Toml_Desc(t *testing.T) {
	bs, _ := Encode(&ownerInfo{
		Org: "org",
	}, nil)
	if !strings.Contains(string(bs), "Org Description") {
		t.Fatal("desc miss")
	}
}

func Test_Toml_Func(t *testing.T) {
	obj := &richStruct{
		Name: "org",
		Driver: func(name string) error {
			return nil
		},
	}
	bs, _ := Encode(obj, nil)
	if !strings.Contains(string(bs), "org") {
		t.Fatalf("struct/func encode fail : %s", string(bs))
	}

	var obj2 richStruct
	if err := Decode(bs, &obj2); err != nil {
		t.Fatalf("struct/func decode fail : %s", err.Error())
	} else if bs2, _ := Encode(obj2, nil); bytes.Compare(bs, bs2) != 0 {
		t.Fatalf("struct/func decode/encode diff fail")
	}
}

func Benchmark_Toml_Decode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var item tomlConfig
		Decode(tomlRaw, &item)
	}
}

func Benchmark_Toml_Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bs, err := Encode(&objRaw, nil)
		if err != nil || len(bs) < 100 {
			//
		}
	}
}

func Benchmark_Json_Decode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var item tomlConfig
		json.Unmarshal(jsonRaw, &item)
	}
}

func Benchmark_Json_Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bs, err := json.Marshal(&objRaw)
		if err != nil || len(bs) < 100 {
			//
		}
	}
}
