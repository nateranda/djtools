package engine

import (
	"bytes"
	"compress/zlib"
	"database/sql"
	"encoding/binary"
	"errors"
	"io"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Library struct {
	songs              []SongNull
	historyList        []HistoryListEntity
	perfData           []PerformanceDataEntry
	playlists          []playlist
	playlistEntityList []playlistEntity
	smartlistList      []smartlist
}

type SongNull struct {
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

type HistoryListEntity struct {
	trackId   int
	startTime int
}

// unused
type SongHistory struct {
	songID     int
	plays      int
	lastPlayed int
}

type PerformanceDataEntry struct {
	trackId   int
	trackData []byte
	beatData  []byte
	quickCues []byte
	loops     []byte
}

type playlist struct {
	playlistId   int
	title        string
	parentListId int
	nextListId   int
}

type playlistEntity struct {
	listId       int
	trackId      int
	nextEntityId int
}

type smartlist struct {
	listUuid           string
	title              string
	parentPlaylistPath sql.NullString
	nextPlaylistPath   sql.NullString
	nextListUuid       sql.NullString
	rules              string
}

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// qUncompress uncompresses a uInt32-appended byte slice using zlib
// used for blobs compressed with the QT C++ library's qCompress function
func qUncompress(file []byte) ([]byte, error) {
	uncompressLength := binary.BigEndian.Uint32(file[:4])
	buffer := bytes.NewBuffer(file[4:])
	r, err := zlib.NewReader(buffer)
	logError(err)

	defer r.Close()

	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	logError(err)

	fileDecomp := out.Bytes()

	// check if the file's uncompressed length matches the header
	if len(fileDecomp) != int(uncompressLength) {
		err := errors.New("db: uncompressed file length does not match length header")
		return []byte{}, err
	} else {
		return fileDecomp, nil
	}

}
