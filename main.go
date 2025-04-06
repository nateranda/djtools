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
	db.ExtractEngineSongs("./databases/engine/")
}
