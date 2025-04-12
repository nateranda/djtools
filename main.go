package main

import (
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
	m, _ := engine.InitDB("databases/engine/")
	perfData := engine.ImportExtractPerformanceData(m)
	//engine.DLBeatData(perfData[8])
	engine.ImportConvertPerformanceData(perfData)
}
