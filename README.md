# htoml4g

hooto toml-based configuration library for golang.

Example:

``` go
package main

import (
	"fmt"

	"github.com/hooto/htoml4g/htoml"
)

type ConfigEntry struct {
	Title     string           `toml:"title"`
	Databases []ConfigDatabase `toml:"databases"`
}

type ConfigDatabase struct {
	Addr    string `toml:"addr"`
	ConnMax int    `toml:"conn_max"`
}

var (
	confRaw = []byte(`
# comment of example
title = "TOML Example"

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
	if err := htoml.Decode(&item, confRaw); err != nil {
		panic(err)
	}

	fmt.Printf("title : %s\n", item.Title)
	fmt.Printf("databases %d :\n", len(item.Databases))
	for i, v := range item.Databases {
		fmt.Printf("  database #%d, addr %s, conn_max %d\n",
			i, v.Addr, v.ConnMax)
	}
}
```

# Dependent Projects or Documents

* TOML (Tom's Obvious, Minimal Language) [http://github.com/toml-lang/toml](http://github.com/toml-lang/toml)
* TOML parser for Golang with reflection [https://github.com/BurntSushi/toml](https://github.com/BurntSushi/toml)

