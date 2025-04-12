package main

import (
	"fmt"
	"log"

	"github.com/nateranda/djtools/engine"
)

// for dev only, will eventually replace with proper bubble up error handling
func logError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	library, err := engine.Import("./databases/engine/")
	logError(err)
	fmt.Printf("song: %v\n", library.Songs[0])
}
