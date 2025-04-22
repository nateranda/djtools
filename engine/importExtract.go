package engine

import (
	"database/sql"
	"fmt"
)

func importExtract(path string) (library, error) {
	var enLibrary library
	var err error

	m, hm, err := initDB(path)
	if err != nil {
		return library{}, err
	}
	enLibrary.songs, err = importExtractTrack(m)
	if err != nil {
		return library{}, fmt.Errorf("error extracting track data: %v", err)
	}
	enLibrary.songHistoryList, err = importExtractHistory(hm)
	if err != nil {
		return library{}, fmt.Errorf("error extracting history data: %v", err)
	}
	enLibrary.perfData, err = importExtractPerformanceData(m)
	if err != nil {
		return library{}, fmt.Errorf("error extracting performance data: %v", err)
	}
	enLibrary.playlists, err = importExtractPlaylist(m)
	if err != nil {
		return library{}, fmt.Errorf("error extracting playlists: %v", err)
	}
	enLibrary.playlistEntityList, err = importExtractPlaylistEntity(m)
	if err != nil {
		return library{}, fmt.Errorf("error extracting playlist data: %v", err)
	}
	enLibrary.smartlistList, err = importExtractSmartlist(m)
	if err != nil {
		return library{}, fmt.Errorf("error extracting smartlists: %v", err)
	}
	return enLibrary, nil
}

func importExtractTrack(db *sql.DB) ([]songNull, error) {
	query := `SELECT id, title, artist, composer, album, genre, fileType, fileBytes, length, year,
		bpm, dateAdded, bitrate, comment, rating, path, remixer, key, label, lastEditTime
		FROM Track ORDER BY id`

	return queryAndScanRows(db, query, func(r *sql.Rows) (songNull, error) {
		var song songNull
		err := r.Scan(
			&song.id, &song.title, &song.artist, &song.composer, &song.album, &song.genre, &song.filetype,
			&song.size, &song.length, &song.year, &song.bpm, &song.dateAdded, &song.bitrate, &song.comment,
			&song.rating, &song.path, &song.remixer, &song.key, &song.label, &song.lastEditTime,
		)
		return song, err
	})
}

func importExtractHistory(db *sql.DB) ([]songHistory, error) {
	query := `SELECT Track.originTrackId, COUNT(HistorylistEntity.trackId), MAX(HistorylistEntity.startTime) 
		FROM Track JOIN HistorylistEntity ON Track.id=HistorylistEntity.trackId
		GROUP BY Track.originTrackId ORDER BY Track.originTrackId`

	return queryAndScanRows(db, query, func(r *sql.Rows) (songHistory, error) {
		var songHistory songHistory
		err := r.Scan(&songHistory.id, &songHistory.plays, &songHistory.lastPlayed)
		return songHistory, err
	})
}

func importExtractPerformanceData(db *sql.DB) ([]performanceDataEntry, error) {
	query := `SELECT trackId, beatData, quickCues, loops FROM PerformanceData ORDER BY trackId`

	return queryAndScanRows(db, query, func(r *sql.Rows) (performanceDataEntry, error) {
		var perfData performanceDataEntry
		err := r.Scan(&perfData.id, &perfData.beatDataBlob, &perfData.quickCuesBlob, &perfData.loopsBlob)
		return perfData, err
	})
}

func importExtractPlaylist(db *sql.DB) ([]playlist, error) {
	query := `SELECT id, title, parentListId, nextListId FROM Playlist ORDER BY id`

	return queryAndScanRows(db, query, func(r *sql.Rows) (playlist, error) {
		var playlist playlist
		err := r.Scan(&playlist.id, &playlist.title, &playlist.parentListId, &playlist.nextListId)
		return playlist, err
	})
}

func importExtractPlaylistEntity(db *sql.DB) ([]playlistEntity, error) {
	query := `SELECT id, listId, trackId, nextEntityId FROM PlaylistEntity ORDER BY listId`

	return queryAndScanRows(db, query, func(r *sql.Rows) (playlistEntity, error) {
		var playlistEntity playlistEntity
		err := r.Scan(&playlistEntity.id, &playlistEntity.listId,
			&playlistEntity.trackId, &playlistEntity.nextEntityId)
		return playlistEntity, err
	})
}

func importExtractSmartlist(db *sql.DB) ([]smartlist, error) {
	query := `SELECT listUuid, title, parentPlaylistPath, nextPlaylistPath, nextListUuid, rules
		FROM Smartlist ORDER BY listUuid`

	return queryAndScanRows(db, query, func(r *sql.Rows) (smartlist, error) {
		var smartlist smartlist
		err := r.Scan(&smartlist.listUuid, &smartlist.title, &smartlist.parentPlaylistPath,
			&smartlist.nextPlaylistPath, &smartlist.nextListUuid, &smartlist.rules)
		return smartlist, err
	})
}

// queryAndScanRows queries a given database and scans
// each row in the response based on a given function.
func queryAndScanRows[T any](db *sql.DB, query string, scanFunc func(*sql.Rows) (T, error)) ([]T, error) {
	r, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query '%s': %v", query, err)
	}
	defer r.Close()

	var results []T
	for r.Next() {
		item, err := scanFunc(r)
		if err != nil {
			return nil, fmt.Errorf("scan error: %v", err)
		}
		results = append(results, item)
	}
	return results, nil
}
