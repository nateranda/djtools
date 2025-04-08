package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type songNull struct {
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

type historyListEntity struct {
	trackId   int
	startTime int
}

type songHistory struct {
	songID     int
	plays      int
	lastPlayed int
}

func songNullCorrect(song songNull) Song {
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

func enImportExtractTrack(db *sql.DB, songs []Song) []Song {
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
		songNull := songNull{}
		err := r.Scan(
			&songNull.SongID,
			&songNull.Title,
			&songNull.Artist,
			&songNull.Composer,
			&songNull.Album,
			&songNull.Genre,
			&songNull.Filetype,
			&songNull.Size,
			&songNull.Length,
			&songNull.Year,
			&songNull.Bpm,
			&songNull.DateAdded,
			&songNull.Bitrate,
			&songNull.Comment,
			&songNull.Rating,
			&songNull.Path,
			&songNull.Remixer,
			&songNull.Key,
			&songNull.Label,
		)
		logError(err)
		song := songNullCorrect(songNull)
		songs = append(songs, song)
	}

	return songs
}

func getSongHistoryData(historyList []historyListEntity) []songHistory {
	var songId int
	var lastPlayed int
	plays := 1

	var songHistoryData []songHistory

	for i, historyListEntity := range historyList {
		if historyListEntity.trackId > songId && i != 0 {
			songHistoryData = append(songHistoryData, songHistory{songId, plays, lastPlayed})
			plays = 0
		}
		songId = historyListEntity.trackId
		lastPlayed = historyListEntity.startTime
		plays += 1
	}
	songHistoryData = append(songHistoryData, songHistory{songId, plays, lastPlayed})

	return songHistoryData
}

func enImportExtractHistory(songs []Song, db_path string) []Song {
	db, err := sql.Open("sqlite3", db_path+"hm.db")
	logError(err)
	defer db.Close()

	query := `SELECT Track.originTrackId, HistorylistEntity.startTime
		FROM Track JOIN HistorylistEntity ON Track.id=HistorylistEntity.trackId
		ORDER BY originTrackId, startTime`

	historyList := []historyListEntity{}

	r, err := db.Query(query)
	logError(err)
	defer r.Close()

	for r.Next() {
		historyListEntity := historyListEntity{}
		startTime := time.Time{}
		err := r.Scan(&historyListEntity.trackId, &startTime)
		logError(err)
		historyListEntity.startTime = int(startTime.Unix())
		historyList = append(historyList, historyListEntity)
	}

	// move this to enImportConvert
	songHistoryData := getSongHistoryData(historyList)
	fmt.Println(songHistoryData)
	return songs
}

func EnImportExtract(library Library, path string) (Library, error) {
	db, err := sql.Open("sqlite3", path+"m.db")
	logError(err)
	defer db.Close()

	var songs []Song

	songs = enImportExtractTrack(db, songs)
	songs = enImportExtractHistory(songs, path)
	fmt.Println(songs[0])
	return Library{}, nil
}

// TBI
func enImportConvert(library Library) (Library, error) {
	return library, nil
}

// TBI
func enImportInject(library Library) (Library, error) {
	return library, nil
}

// TBI
func EnImport(path string, options ImportOptions) (Library, error) {
	var library Library

	library, err := EnImportExtract(library, path)
	logError(err)
	library, err = enImportConvert(library)
	logError(err)
	library, err = enImportInject(library)
	logError(err)

	return library, nil
}

// TBI
func EnExport(library Library, path string, options ExportOptions) error {
	return nil
}
