package main

import (
	"log"

	"github.com/nateranda/djtools/db"
	"github.com/nateranda/djtools/engine"
)

func main() {
	importOptions := engine.ImportOptions{}
	library, err := engine.Import("/Users/nateranda/Music/Engine Library Test/", importOptions)
	if err != nil {
		log.Fatal(err)
	}
	db.Save(&library, "./tmp/library")
	//var library db.Library

	//db.Load(&library, "./tmp/library")
	//rbxml.Export(&library, "./tmp/library.xml")
}
