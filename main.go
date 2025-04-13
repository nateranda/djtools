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
	library, err := engine.Import("/Users/nateranda/Music/Engine Library/Database2/")
	logError(err)
	for _, song := range library.Playlists[2].Songs {
		fmt.Println(song.Title)
	}
}
