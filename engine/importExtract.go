package engine

import (
	"database/sql"
	"fmt"
)

// queryAndScanRows is a generic helper that queries a given database and scans
// each row in the response based on a given function.
func queryAndScanRows[T any](db *sql.DB, query string, scanFunc func(*sql.Rows) (T, error)) ([]T, error) {
	r, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("query error: %v", err)
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

func importExtractTrack(engineDB engineDB) ([]songNull, error) {
	query := `SELECT id, title, artist, composer, album, genre, fileType, fileBytes, length, year,
		bpm, dateAdded, bitrate, comment, rating, path, remixer, key, label, lastEditTime
		FROM Track ORDER BY id`

	return queryAndScanRows(engineDB.m, query, func(r *sql.Rows) (songNull, error) {
		var song songNull
		err := r.Scan(
			&song.id, &song.title, &song.artist, &song.composer, &song.album, &song.genre, &song.filetype,
			&song.size, &song.length, &song.year, &song.bpm, &song.dateAdded, &song.bitrate, &song.comment,
			&song.rating, &song.path, &song.remixer, &song.key, &song.label, &song.lastEditTime,
		)
		return song, err
	})
}

func importExtractHistory(engineDB engineDB) ([]songHistory, error) {
	query := `SELECT Track.originTrackId, COUNT(HistorylistEntity.trackId), MAX(HistorylistEntity.startTime) 
		FROM Track JOIN HistorylistEntity ON Track.id=HistorylistEntity.trackId
		GROUP BY Track.originTrackId ORDER BY Track.originTrackId`

	return queryAndScanRows(engineDB.hm, query, func(r *sql.Rows) (songHistory, error) {
		var songHistory songHistory
		err := r.Scan(&songHistory.id, &songHistory.plays, &songHistory.lastPlayed)
		return songHistory, err
	})
}

func importExtractPerformanceData(engineDB engineDB) ([]performanceDataEntry, error) {
	query := `SELECT trackId, beatData, quickCues, loops FROM PerformanceData ORDER BY trackId`

	return queryAndScanRows(engineDB.m, query, func(r *sql.Rows) (performanceDataEntry, error) {
		var perfData performanceDataEntry
		err := r.Scan(&perfData.id, &perfData.beatDataBlob, &perfData.quickCuesBlob, &perfData.loopsBlob)
		return perfData, err
	})
}

func importExtractPlaylist(engineDB engineDB) ([]playlist, error) {
	query := `SELECT id, title, parentListId, nextListId FROM Playlist ORDER BY id`

	return queryAndScanRows(engineDB.m, query, func(r *sql.Rows) (playlist, error) {
		var playlist playlist
		err := r.Scan(&playlist.id, &playlist.title, &playlist.parentListId, &playlist.nextListId)
		return playlist, err
	})
}

func importExtractPlaylistEntity(engineDB engineDB) ([]playlistEntity, error) {
	query := `SELECT id, listId, trackId, nextEntityId FROM PlaylistEntity ORDER BY listId`

	return queryAndScanRows(engineDB.m, query, func(r *sql.Rows) (playlistEntity, error) {
		var playlistEntity playlistEntity
		err := r.Scan(&playlistEntity.id, &playlistEntity.listId,
			&playlistEntity.trackId, &playlistEntity.nextEntityId)
		return playlistEntity, err
	})
}

func importExtractSmartlist(engineDB engineDB) ([]smartlist, error) {
	query := `SELECT listUuid, title, parentPlaylistPath, nextPlaylistPath, nextListUuid, rules FROM Smartlist ORDER BY listUuid`

	return queryAndScanRows(engineDB.m, query, func(r *sql.Rows) (smartlist, error) {
		var smartlist smartlist
		err := r.Scan(&smartlist.listUuid, &smartlist.title, &smartlist.parentPlaylistPath,
			&smartlist.nextPlaylistPath, &smartlist.nextListUuid, &smartlist.rules)
		return smartlist, err
	})
}

func importExtract(path string) (library, error) {
	var enLibrary library
	var err error

	engineDB, err := initDB(path)
	if err != nil {
		return library{}, err
	}
	enLibrary.songs, err = importExtractTrack(engineDB)
	if err != nil {
		return library{}, err
	}
	enLibrary.songHistoryList, err = importExtractHistory(engineDB)
	if err != nil {
		return library{}, err
	}
	enLibrary.perfData, err = importExtractPerformanceData(engineDB)
	if err != nil {
		return library{}, err
	}
	enLibrary.playlists, err = importExtractPlaylist(engineDB)
	if err != nil {
		return library{}, err
	}
	enLibrary.playlistEntityList, err = importExtractPlaylistEntity(engineDB)
	if err != nil {
		return library{}, err
	}
	enLibrary.smartlistList, err = importExtractSmartlist(engineDB)
	if err != nil {
		return library{}, err
	}
	return enLibrary, nil
}
