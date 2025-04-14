package main

import (
	"fmt"
	"log"

	"github.com/nateranda/djtools/engine"
)

func main() {
	library, err := engine.Import("/Users/nateranda/Music/Engine Library/Database2/")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", library)
}
