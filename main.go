package main

import (
	"log"

	"github.com/nateranda/djtools/db"
	"github.com/nateranda/djtools/engine"
	"github.com/nateranda/djtools/rbxml"
)

func main() {
	importOptions := engine.ImportOptions{}
	library, err := engine.Import("/Users/nateranda/Music/Engine Library/", importOptions)
	if err != nil {
		log.Fatal(err)
	}
	db.Save(&library, "./tmp/library")
	//var library db.Library

	db.Load(&library, "./tmp/library")
	rbxml.Export(&library, "./tmp/library.xml")
}
