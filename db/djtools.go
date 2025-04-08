package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitLibrary() {
	db, err := sql.Open("sqlite3", "./library.db")
	logError(err)

	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS songs (
		songID INTEGER PRIMARY KEY,
		title TEXT,
		artist TEXT,
		composer TEXT,
		album TEXT,
		grouping TEXT,
		genre TEXT,
		type TEXT,
		size INT,
		length REAL,
		albumArt BLOB,
		trackNumber INT,
		year INT,
		bpm REAL,
		dateModified INT,
		bitrate INT,
		sampleRate REAL,
		comment TEXT,
		playCount INT,
		lastPlayed INT,
		rating INT,
		path TEXT,
		remixer TEXT,
		key TEXT,
		label TEXT,
		mix TEXT,
		color TEXT
	)`)
	logError(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS grids (
		startPosition REAL,
		bpm REAL,
		songID INT,
		beatNumber INT
	)`)
	logError(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cues (
		name TEXT,
		songID INT,
		offset REAL,
		position INT
	)`)
	logError(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS loops (
		name TEXT,
		songID INT,
		offset REAL,
		position INT,
		length REAL
	)`)
	logError(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS playlists (
		playlistID INT PRIMARY KEY,
		position INT,
		name TEXT,
		smartPlaylist INT,
		parameters TEXT,
		folderID INT
	)`)
	logError(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS playlistContent (
		songID INT,
		playlistID INT,
		position INT
	)`)
	logError(err)

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS folders (
		folderID INT PRIMARY KEY,
		position INT,
		name TEXT
	)`)
	logError(err)
}
