package engine

import (
	"bytes"
	"compress/zlib"
	"database/sql"
	"encoding/binary"
	"errors"
	"io"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type library struct {
	songs              []songNull
	historyList        []historyListEntity
	perfData           []performanceDataEntry
	playlists          []playlist
	playlistEntityList []playlistEntity
	smartlistList      []smartlist
}

type songNull struct {
	id        sql.NullInt64
	title     sql.NullString
	artist    sql.NullString
	composer  sql.NullString
	album     sql.NullString
	genre     sql.NullString
	filetype  sql.NullString
	size      sql.NullInt64
	length    sql.NullFloat64
	year      sql.NullInt64
	bpm       sql.NullFloat64
	dateAdded sql.NullTime
	bitrate   sql.NullInt64
	comment   sql.NullString
	rating    sql.NullInt64
	path      sql.NullString
	remixer   sql.NullString
	key       sql.NullString
	label     sql.NullString
}

type historyListEntity struct {
	id        int
	startTime time.Time
}

// unused
type songHistory struct {
	id         int
	plays      int
	lastPlayed int
}

type performanceDataEntry struct {
	id        int
	trackData []byte
	beatData  []byte
	quickCues []byte
	loops     []byte
}

type playlist struct {
	id           int
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

type beatData struct {
	sampleRate      float64
	defaultBeatgrid []marker
	adjBeatgrid     []marker
}

type marker struct {
	offset     float64
	beatNumber int64
	numBeats   uint32
}

func logError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// qUncompress uncompresses a uInt32-appended byte slice using zlib,
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

// InitDB initializes the Engine SQL database at a given path
func InitDB(path string) (*sql.DB, *sql.DB) {
	m, err := sql.Open("sqlite3", path+"m.db")
	logError(err)
	hm, err := sql.Open("sqlite3", path+"hm.db")
	logError(err)

	return m, hm
}
