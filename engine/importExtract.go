package engine

import (
	"database/sql"
	"time"
)

func importExtractTrack(db *sql.DB) []SongNull {
	var songs []SongNull

	query := `SELECT id, title, artist, composer,
		album, genre, fileType, fileBytes,
		length, year, bpm, dateAdded,
		bitrate, comment, rating, path,
		remixer, key, label
		FROM Track ORDER BY id`

	r, err := db.Query(query)
	logError(err)
	defer r.Close()

	for r.Next() {
		song := SongNull{}
		err := r.Scan(
			&song.SongID,
			&song.Title,
			&song.Artist,
			&song.Composer,
			&song.Album,
			&song.Genre,
			&song.Filetype,
			&song.Size,
			&song.Length,
			&song.Year,
			&song.Bpm,
			&song.DateAdded,
			&song.Bitrate,
			&song.Comment,
			&song.Rating,
			&song.Path,
			&song.Remixer,
			&song.Key,
			&song.Label,
		)
		logError(err)
		songs = append(songs, song)
	}

	return songs
}

func importExtractHistory(db *sql.DB) []HistoryListEntity {
	query := `SELECT Track.originTrackId, HistorylistEntity.startTime
		FROM Track JOIN HistorylistEntity ON Track.id=HistorylistEntity.trackId
		ORDER BY originTrackId, startTime`

	var historyList []HistoryListEntity

	r, err := db.Query(query)
	logError(err)
	defer r.Close()

	for r.Next() {
		HistoryListEntity := HistoryListEntity{}
		startTime := time.Time{}
		err := r.Scan(&HistoryListEntity.trackId, &startTime)
		logError(err)
		HistoryListEntity.startTime = int(startTime.Unix())
		historyList = append(historyList, HistoryListEntity)
	}

	return historyList
}

func importExtractPerformanceData(db *sql.DB) {
	query := `SELECT trackId, trackData, beatData, quickCues, loops FROM PerformanceData ORDER BY trackId`

	var perfDataList []PerformanceDataEntry

	r, err := db.Query(query)
	logError(err)
	defer r.Close()

	for r.Next() {
		var perfData PerformanceDataEntry
		err := r.Scan(&perfData.trackId, &perfData.trackData, &perfData.beatData, &perfData.quickCues, &perfData.loops)
		logError(err)
		perfDataList = append(perfDataList, perfData)
	}

	DLBeatData(perfDataList[4])
}
