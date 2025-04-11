package engine

import (
	"os"

	"github.com/nateranda/djtools/db"
)

func DLBeatData(perfData performanceDataEntry) {
	err := os.WriteFile("tmp/trackData", perfData.trackData, 0644)
	logError(err)
	err = os.WriteFile("tmp/beatData", perfData.beatData, 0644)
	logError(err)
	err = os.WriteFile("tmp/quickCues", perfData.quickCues, 0644)
	logError(err)
	err = os.WriteFile("tmp/loops", perfData.loops, 0644)
	logError(err)
}

func ImportExtract(path string) (library, error) {
	var Library library

	m, hm := InitDB(path)
	Library.songs = importExtractTrack(m)
	Library.historyList = importExtractHistory(hm)
	Library.perfData = ImportExtractPerformanceData(m)
	Library.playlists = importExtractPlaylist(m)
	Library.playlistEntityList = importExtractPlaylistEntity(m)
	Library.smartlistList = importExtractSmartlist(m)

	return Library, nil
}

func ImportConvert(enLibrary library) (db.Library, error) {
	var library db.Library
	library.Songs = importConvertSongList(enLibrary.songs)
	return library, nil
}
