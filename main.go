package main

import (
	"fmt"
	"log"

	"github.com/nateranda/djtools/engine"
)

func main() {
	importOptions := engine.ImportOptions{}

	library, err := engine.Import("/Users/nateranda/Music/Engine Library/Database2/", importOptions)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v", library.Songs[1].Grid)
}
