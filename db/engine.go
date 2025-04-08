package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type enLibrary struct {
	songs       []enSongNull
	historyList []enHistoryListEntity
}

type enSongNull struct {
	SongID    sql.NullInt64
	Title     sql.NullString
	Artist    sql.NullString
	Composer  sql.NullString
	Album     sql.NullString
	Genre     sql.NullString
	Filetype  sql.NullString
	Size      sql.NullInt64
	Length    sql.NullFloat64
	Year      sql.NullInt64
	Bpm       sql.NullFloat64
	DateAdded sql.NullTime
	Bitrate   sql.NullInt64
	Comment   sql.NullString
	Rating    sql.NullInt64
	Path      sql.NullString
	Remixer   sql.NullString
	Key       sql.NullString
	Label     sql.NullString
}

type enHistoryListEntity struct {
	trackId   int
	startTime int
}

// unused
type enSongHistory struct {
	songID     int
	plays      int
	lastPlayed int
}

func enImportExtractTrack(db *sql.DB) []enSongNull {
	var songs []enSongNull

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
		song := enSongNull{}
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

func enImportExtractHistory(db *sql.DB) []enHistoryListEntity {
	query := `SELECT Track.originTrackId, HistorylistEntity.startTime
		FROM Track JOIN HistorylistEntity ON Track.id=HistorylistEntity.trackId
		ORDER BY originTrackId, startTime`

	var historyList []enHistoryListEntity

	r, err := db.Query(query)
	logError(err)
	defer r.Close()

	for r.Next() {
		enHistoryListEntity := enHistoryListEntity{}
		startTime := time.Time{}
		err := r.Scan(&enHistoryListEntity.trackId, &startTime)
		logError(err)
		enHistoryListEntity.startTime = int(startTime.Unix())
		historyList = append(historyList, enHistoryListEntity)
	}

	return historyList
}

func EnImportExtract(path string) (enLibrary, error) {
	var enLibrary enLibrary

	dbm, err := sql.Open("sqlite3", path+"m.db")
	logError(err)
	defer dbm.Close()

	dbhm, err := sql.Open("sqlite3", path+"hm.db")
	logError(err)
	defer dbhm.Close()

	enLibrary.songs = enImportExtractTrack(dbm)
	enLibrary.historyList = enImportExtractHistory(dbhm)

	fmt.Println(enLibrary.songs[0])
	fmt.Println(enLibrary.historyList[0])
	return enLibrary, nil
}

// unused
func enImportConvertSong(song enSongNull) Song {
	return Song{
		SongID:    int(song.SongID.Int64),
		Title:     song.Title.String,
		Artist:    song.Artist.String,
		Composer:  song.Composer.String,
		Album:     song.Album.String,
		Genre:     song.Genre.String,
		Filetype:  song.Filetype.String,
		Size:      int(song.Size.Int64),
		Length:    float32(song.Length.Float64),
		Year:      int(song.Year.Int64),
		Bpm:       float32(song.Bpm.Float64),
		DateAdded: int(song.DateAdded.Time.Unix()),
		Bitrate:   int(song.Bitrate.Int64),
		Comment:   song.Comment.String,
		Rating:    int(song.Rating.Int64),
		Path:      song.Path.String,
		Remixer:   song.Remixer.String,
		Key:       song.Key.String,
		Label:     song.Label.String,
	}
}

// unused
func enImportConvertSongHistory(historyList []enHistoryListEntity) []enSongHistory {
	var songId int
	var lastPlayed int
	plays := 1

	var enSongHistoryData []enSongHistory

	for i, enHistoryListEntity := range historyList {
		if enHistoryListEntity.trackId > songId && i != 0 {
			enSongHistoryData = append(enSongHistoryData, enSongHistory{songId, plays, lastPlayed})
			plays = 0
		}
		songId = enHistoryListEntity.trackId
		lastPlayed = enHistoryListEntity.startTime
		plays += 1
	}
	enSongHistoryData = append(enSongHistoryData, enSongHistory{songId, plays, lastPlayed})

	return enSongHistoryData
}

// TBI
func enImportConvert(enLibrary enLibrary) (Library, error) {
	return Library{}, nil
}

// TBI
func enImportInject(library Library) (Library, error) {
	return library, nil
}

// TBI
func EnImport(path string, options ImportOptions) (Library, error) {
	var library Library

	enLibrary, err := EnImportExtract(path)
	logError(err)
	library, err = enImportConvert(enLibrary)
	logError(err)
	library, err = enImportInject(library)
	logError(err)

	return library, nil
}

// TBI
func EnExport(library Library, path string, options ExportOptions) error {
	return nil
}
