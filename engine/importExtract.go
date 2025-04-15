package engine

import (
	"database/sql"
	"fmt"
)

func importExtractTrack(db *sql.DB) ([]songNull, error) {
	var songs []songNull

	query := `SELECT id, title, artist, composer,
		album, genre, fileType, fileBytes,
		length, year, bpm, dateAdded,
		bitrate, comment, rating, path,
		remixer, key, label, lastEditTime
		FROM Track ORDER BY id`

	r, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error extracting tracks: %v", err)
	}
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
			&song.lastEditTime,
		)
		if err != nil {
			return nil, fmt.Errorf("error extracting tracks: %v", err)
		}
		songs = append(songs, song)
	}

	return songs, nil
}

func importExtractHistory(db *sql.DB) ([]historyListEntity, error) {
	query := `SELECT Track.originTrackId, HistorylistEntity.startTime
		FROM Track JOIN HistorylistEntity ON Track.id=HistorylistEntity.trackId
		ORDER BY originTrackId, startTime`

	var historyList []historyListEntity

	r, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error extracting track history: %v", err)
	}
	defer r.Close()

	for r.Next() {
		HistoryListEntity := historyListEntity{}
		err := r.Scan(&HistoryListEntity.id, &HistoryListEntity.startTime)
		if err != nil {
			return nil, fmt.Errorf("error extracting track history: %v", err)
		}
		historyList = append(historyList, HistoryListEntity)
	}

	return historyList, nil
}

func importExtractPerformanceData(db *sql.DB) ([]performanceDataEntry, error) {
	query := `SELECT trackId, beatData, quickCues, loops FROM PerformanceData ORDER BY trackId`

	var perfDataList []performanceDataEntry

	r, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error extracting performance data: %v", err)
	}
	defer r.Close()

	for r.Next() {
		var perfData performanceDataEntry
		err := r.Scan(&perfData.id, &perfData.beatDataBlob, &perfData.quickCuesBlob, &perfData.loopsBlob)
		if err != nil {
			return nil, fmt.Errorf("error extracting performance data: %v", err)
		}
		perfDataList = append(perfDataList, perfData)
	}

	return perfDataList, nil
}

func importExtractPlaylist(db *sql.DB) ([]playlist, error) {
	query := `SELECT id, title, parentListId, nextListId FROM Playlist ORDER BY id`

	var playlists []playlist

	r, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error extracting playlists: %v", err)
	}
	defer r.Close()

	for r.Next() {
		var playlist playlist
		err := r.Scan(&playlist.id, &playlist.title, &playlist.parentListId, &playlist.nextListId)
		if err != nil {
			return nil, fmt.Errorf("error extracting playlists: %v", err)
		}
		playlists = append(playlists, playlist)
	}

	return playlists, nil
}

func importExtractPlaylistEntity(db *sql.DB) ([]playlistEntity, error) {
	query := `SELECT id, listId, trackId, nextEntityId FROM PlaylistEntity ORDER BY listId`

	var playlistEntityList []playlistEntity

	r, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error extracting playlist entities: %v", err)
	}
	defer r.Close()

	for r.Next() {
		var playlistEntity playlistEntity
		err := r.Scan(&playlistEntity.id, &playlistEntity.listId, &playlistEntity.trackId, &playlistEntity.nextEntityId)
		if err != nil {
			return nil, fmt.Errorf("error extracting playlist entities: %v", err)
		}
		playlistEntityList = append(playlistEntityList, playlistEntity)
	}

	return playlistEntityList, nil
}

func importExtractSmartlist(db *sql.DB) ([]smartlist, error) {
	query := `SELECT listUuid, title, parentPlaylistPath, nextPlaylistPath, nextListUuid, rules FROM Smartlist ORDER BY listUuid`

	var smartlistList []smartlist

	r, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error extracting smartlists: %v", err)
	}
	defer r.Close()

	for r.Next() {
		var smartlist smartlist
		err := r.Scan(&smartlist.listUuid, &smartlist.title, &smartlist.parentPlaylistPath, &smartlist.nextPlaylistPath, &smartlist.nextListUuid, &smartlist.rules)
		if err != nil {
			return nil, fmt.Errorf("error extracting smartlists: %v", err)
		}
		smartlistList = append(smartlistList, smartlist)
	}

	return smartlistList, nil
}

func importExtract(path string) (library, error) {
	var enLibrary library
	var err error

	m, hm, err := initDB(path)
	if err != nil {
		return library{}, err
	}
	enLibrary.songs, err = importExtractTrack(m)
	if err != nil {
		return library{}, err
	}
	enLibrary.historyList, err = importExtractHistory(hm)
	if err != nil {
		return library{}, err
	}
	enLibrary.perfData, err = importExtractPerformanceData(m)
	if err != nil {
		return library{}, err
	}
	enLibrary.playlists, err = importExtractPlaylist(m)
	if err != nil {
		return library{}, err
	}
	enLibrary.playlistEntityList, err = importExtractPlaylistEntity(m)
	if err != nil {
		return library{}, err
	}
	enLibrary.smartlistList, err = importExtractSmartlist(m)
	if err != nil {
		return library{}, err
	}
	return enLibrary, nil
}
