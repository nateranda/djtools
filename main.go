package main

import (
	"fmt"
	"log"

	"github.com/nateranda/djtools/db"
)

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	library, err := db.EnImportExtract(db.Library{}, "./databases/engine/")
	logError(err)
	fmt.Println(library)
}
