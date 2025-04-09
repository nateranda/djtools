package main

import (
	"log"

	"github.com/nateranda/djtools/db"
)

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	//_, err := db.EnImportExtract("./databases/engine/")
	//logError(err)

	db.EnImportConvertGrid()
}
