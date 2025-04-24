package main

import (
	"fmt"
	"log"

	"github.com/nateranda/djtools/serato"
)

func main() {
	crate, err := serato.ExtractCrate("tmp/Subcrates/crate 2%%crate 3.crate")
	if err != nil {
		log.Panic(err)
	}
	fmt.Printf("crate: %v\n", crate)
}
