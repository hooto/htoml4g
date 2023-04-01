package main

import (
	"fmt"

	"github.com/hooto/htoml4g/htoml"
)

type ConfigEntry struct {
	Title     string           `toml:"title"`
	Databases []ConfigDatabase `toml:"databases"`
	Attrs     []string         `toml:"attrs"`
}

type ConfigDatabase struct {
	Addr    string `toml:"addr"`
	ConnMax int    `toml:"conn_max"`
}

var (
	confRaw = []byte(`
# comment of example
title = "TOML Example"
attrs = ["0000", "1111", "2222", "3333", "4444", "5555", "6666", "7777", "8888", "9999"]

[[databases]]
addr = "192.168.0.1:3001"
conn_max = 100

[[databases]]
addr = "192.168.0.2:3001"
conn_max = 100
`)
)

func main() {

	var item ConfigEntry
	if err := htoml.Decode(confRaw, &item); err != nil {
		panic(err)
	}

	fmt.Printf("title : %s\n", item.Title)
	fmt.Printf("databases %d :\n", len(item.Databases))
	for i, v := range item.Databases {
		fmt.Printf("  database #%d, addr %s, conn_max %d\n",
			i, v.Addr, v.ConnMax)
	}

	if b, err := htoml.Encode(item); err == nil {
		fmt.Println("encode", string(b))
	}
}
