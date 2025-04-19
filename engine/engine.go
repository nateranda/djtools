package engine

import (
	"bytes"
	"compress/zlib"
	"database/sql"
	"encoding/binary"
	"errors"
	"fmt"
	"io"

	_ "github.com/mattn/go-sqlite3"
	"github.com/nateranda/djtools/db"
)

// ImportOptions contains the options used when importing an Engine library.
type ImportOptions struct {
	ImportOriginalGrids   bool
	ImportOriginalCues    bool
	PreserveOriginalPaths bool
}

type library struct {
	songs              []songNull
	songHistoryList    []songHistory
	perfData           []performanceDataEntry
	playlists          []playlist
	playlistEntityList []playlistEntity
	smartlistList      []smartlist
}

type songNull struct {
	id           sql.NullInt64
	title        sql.NullString
	artist       sql.NullString
	composer     sql.NullString
	album        sql.NullString
	genre        sql.NullString
	filetype     sql.NullString
	size         sql.NullInt64
	length       sql.NullFloat64
	year         sql.NullInt64
	bpm          sql.NullFloat64
	dateAdded    sql.NullTime
	bitrate      sql.NullInt64
	comment      sql.NullString
	rating       sql.NullInt64
	path         sql.NullString
	remixer      sql.NullString
	key          sql.NullString
	label        sql.NullString
	lastEditTime sql.NullTime
}

type songHistory struct {
	id         int
	plays      int
	lastPlayed int
}

type performanceDataEntry struct {
	id            int
	beatDataBlob  []byte
	quickCuesBlob []byte
	loopsBlob     []byte
}

type playlist struct {
	id           int
	title        string
	parentListId int
	nextListId   int
	songs        []int
}

type playlistEntity struct {
	id           int
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

type cueData struct {
	cues        []db.HotCue
	cueOriginal float64
	cueModified float64
}

type marker struct {
	offset     float64
	beatNumber int64
	numBeats   uint32
}

type engineDB struct {
	m  *sql.DB
	hm *sql.DB
}

// qUncompress uncompresses a uInt32-appended byte slice using zlib,
// used for blobs compressed with the QT C++ library's qCompress function.
func qUncompress(file []byte) ([]byte, error) {
	if len(file) < 5 {
		return nil, fmt.Errorf("error uncompressing file: blob must contain 5 or more bytes")
	}
	uncompressLength := binary.BigEndian.Uint32(file[:4])
	buffer := bytes.NewBuffer(file[4:])
	r, err := zlib.NewReader(buffer)
	if err != nil {
		return nil, fmt.Errorf("error uncompressing file: %v", err)
	}

	defer r.Close()

	var out bytes.Buffer
	io.Copy(&out, r)

	fileDecomp := out.Bytes()

	// check if the file's uncompressed length matches the header
	if len(fileDecomp) != int(uncompressLength) {
		return []byte{}, errors.New("VerificationError: uncompressed file length does not match length header")
	} else {
		return fileDecomp, nil
	}
}

// initDB initializes the Engine SQL database at a given path.
func initDB(path string) (engineDB, error) {
	var db engineDB
	var err error
	db.m, _ = sql.Open("sqlite3", path+"Database2/m.db")
	err = db.m.Ping()
	if err != nil {
		return engineDB{}, fmt.Errorf("error initializing database: %v", err)
	}
	db.hm, _ = sql.Open("sqlite3", path+"Database2/hm.db")
	err = db.hm.Ping()
	if err != nil {
		return engineDB{}, fmt.Errorf("error initializing database: %v", err)
	}
	return db, nil
}

// Import converts an Engine database into a djtools Library struct
func Import(path string, importOptions ImportOptions) (db.Library, error) {
	enLibrary, err := importExtract(path)
	if err != nil {
		return db.Library{}, err
	}
	library, err := importConvert(enLibrary, path, importOptions)
	if err != nil {
		return db.Library{}, err
	}
	return library, nil
}
