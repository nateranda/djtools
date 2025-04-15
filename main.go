package main

import (
	"fmt"

	"github.com/nateranda/djtools/db"
)

func main() {
	//importOptions := engine.ImportOptions{}
	//library, err := engine.Import("/Users/nateranda/Music/Engine Library/Database2/", importOptions)
	//if err != nil {
	//	log.Fatal(err)
	//}
	var library db.Library

	db.Load(&library, "./tmp/library")
	fmt.Println(library)
}
