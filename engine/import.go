package engine

import (
	"database/sql"
	"os"

	"github.com/nateranda/djtools/db"
)

func DLBeatData(perfData PerformanceDataEntry) {
	err := os.WriteFile("tmp/trackData", perfData.trackData, 0644)
	logError(err)
	err = os.WriteFile("tmp/beatData", perfData.beatData, 0644)
	logError(err)
	err = os.WriteFile("tmp/quickCues", perfData.quickCues, 0644)
	logError(err)
	err = os.WriteFile("tmp/loops", perfData.loops, 0644)
	logError(err)
}

func ImportExtract(path string) (Library, error) {
	var Library Library

	dbm, err := sql.Open("sqlite3", path+"m.db")
	logError(err)
	defer dbm.Close()

	dbhm, err := sql.Open("sqlite3", path+"hm.db")
	logError(err)
	defer dbhm.Close()

	Library.songs = importExtractTrack(dbm)
	Library.historyList = importExtractHistory(dbhm)
	// ImportExtractPerformanceData(dbm)

	return Library, nil
}

// TBI
func importConvert(library Library) (db.Library, error) {
	return db.Library{}, nil
}

// TBI
func ImportInject(library db.Library) (db.Library, error) {
	return library, nil
}

// TBI
func Import(path string, options db.ImportOptions) (db.Library, error) {
	var library db.Library

	engineLibrary, err := ImportExtract(path)
	logError(err)
	library, err = importConvert(engineLibrary)
	logError(err)
	library, err = ImportInject(library)
	logError(err)

	return library, nil
}
