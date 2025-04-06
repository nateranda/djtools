package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nateranda/djtools/lib"
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

func songNullCorrect(song songNull) lib.Song {
	return lib.Song{
		SongID:    int(song.SongID.Int64),
		Title:     song.Title.String,
		Artist:    song.Artist.String,
		Composer:  song.Composer.String,
		Album:     song.Album.String,
		Genre:     song.Genre.String,
		Filetype:  song.Filetype.String,
		Size:      int(song.Size.Int64),
		Length:    float32(song.Length.Float64),
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

func engineTrack(db *sql.DB, songs []lib.Song) []lib.Song {
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

func ExtractEngineSongs(db_path string) {
	db, err := sql.Open("sqlite3", db_path+"m.db")
	logError(err)
	defer db.Close()

	var songs []lib.Song

	songs = engineTrack(db, songs)
	fmt.Println(songs[1])
}
