package engine

import (
	"database/sql"
)

func importExtractTrack(db *sql.DB) []songNull {
	var songs []songNull

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
		song := songNull{}
		err := r.Scan(
			&song.id,
			&song.title,
			&song.artist,
			&song.composer,
			&song.album,
			&song.genre,
			&song.filetype,
			&song.size,
			&song.length,
			&song.year,
			&song.bpm,
			&song.dateAdded,
			&song.bitrate,
			&song.comment,
			&song.rating,
			&song.path,
			&song.remixer,
			&song.key,
			&song.label,
		)
		logError(err)
		songs = append(songs, song)
	}

	return songs
}

func importExtractHistory(db *sql.DB) []historyListEntity {
	query := `SELECT Track.originTrackId, HistorylistEntity.startTime
		FROM Track JOIN HistorylistEntity ON Track.id=HistorylistEntity.trackId
		ORDER BY originTrackId, startTime`

	var historyList []historyListEntity

	r, err := db.Query(query)
	logError(err)
	defer r.Close()

	for r.Next() {
		HistoryListEntity := historyListEntity{}
		err := r.Scan(&HistoryListEntity.id, &HistoryListEntity.startTime)
		logError(err)
		historyList = append(historyList, HistoryListEntity)
	}

	return historyList
}

func ImportExtractPerformanceData(db *sql.DB) []performanceDataEntry {
	query := `SELECT trackId, beatData, quickCues, loops FROM PerformanceData ORDER BY trackId`

	var perfDataList []performanceDataEntry

	r, err := db.Query(query)
	logError(err)
	defer r.Close()

	for r.Next() {
		var perfData performanceDataEntry
		err := r.Scan(&perfData.id, &perfData.beatDataBlob, &perfData.quickCuesBlob, &perfData.loopsBlob)
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
		err := r.Scan(&playlist.id, &playlist.title, &playlist.parentListId, &playlist.nextListId)
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
