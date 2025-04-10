package main

import (
	"log"
)

func logError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {}
