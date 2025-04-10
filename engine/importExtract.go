package engine

import (
	"database/sql"
)

func ImportExtractInitDB(path string) (*sql.DB, *sql.DB) {
	m, err := sql.Open("sqlite3", path+"m.db")
	logError(err)
	hm, err := sql.Open("sqlite3", path+"hm.db")
	logError(err)

	return m, hm
}

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
		err := r.Scan(&HistoryListEntity.trackId, &HistoryListEntity.startTime)
		logError(err)
		historyList = append(historyList, HistoryListEntity)
	}

	return historyList
}

func importExtractPerformanceData(db *sql.DB) []PerformanceDataEntry {
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

	return perfDataList
}

func importExtractPlaylist(db *sql.DB) []playlist {
	query := `SELECT id, title, parentListId, nextListId FROM Playlist ORDER BY id`

	var playlists []playlist

	r, err := db.Query(query)
	logError(err)
	defer r.Close()

	for r.Next() {
		var playlist playlist
		err := r.Scan(&playlist.playlistId, &playlist.title, &playlist.parentListId, &playlist.nextListId)
		logError(err)
		playlists = append(playlists, playlist)
	}

	return playlists
}

func importExtractPlaylistEntity(db *sql.DB) []playlistEntity {
	query := `SELECT listId, trackId, nextEntityId FROM PlaylistEntity ORDER BY trackId`

	var playlistEntityList []playlistEntity

	r, err := db.Query(query)
	logError(err)
	defer r.Close()

	for r.Next() {
		var playlistEntity playlistEntity
		err := r.Scan(&playlistEntity.listId, &playlistEntity.trackId, &playlistEntity.nextEntityId)
		logError(err)
		playlistEntityList = append(playlistEntityList, playlistEntity)
	}

	return playlistEntityList
}

func importExtractSmartlist(db *sql.DB) []smartlist {
	query := `SELECT listUuid, title, parentPlaylistPath, nextPlaylistPath, nextListUuid, rules FROM Smartlist ORDER BY listUuid`

	var smartlistList []smartlist

	r, err := db.Query(query)
	logError(err)
	defer r.Close()

	for r.Next() {
		var smartlist smartlist
		err := r.Scan(&smartlist.listUuid, &smartlist.title, &smartlist.parentPlaylistPath, &smartlist.nextPlaylistPath, &smartlist.nextListUuid, &smartlist.rules)
		logError(err)
		smartlistList = append(smartlistList, smartlist)
	}

	return smartlistList
}
