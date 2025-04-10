package main

import (
	"log"
)

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {}
