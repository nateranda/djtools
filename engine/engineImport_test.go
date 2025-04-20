package engine_test

import (
	"database/sql"
	"encoding/json"
	"errors"
	"os"
	"sort"
	"testing"

	"github.com/nateranda/djtools/db"
	"github.com/nateranda/djtools/engine"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

// preserve original relative filepaths for operating system parity
var defaultOptions = engine.ImportOptions{PreserveOriginalPaths: true}

// generateDatabase generates an Engine database from m.sql and hm.sql files
func generateDatabase(t *testing.T, fixturePath string) string {
	tempdir := t.TempDir() + "/"

	//make Database2 directory inside of temp directory
	os.Mkdir(tempdir+"Database2/", 0755)

	// open and populate m.db with given fixture
	m, _ := sql.Open("sqlite3", tempdir+"Database2/m.db")
	err := m.Ping()
	if err != nil {
		t.Errorf("unexpected error creating test database: %v", err)
	}

	queryByte, err := os.ReadFile(fixturePath + "m.sql")
	if err != nil {
		t.Errorf("unexpected error reading from m.db fixture: %v", err)
	}
	query := string(queryByte)

	m.Exec(query)

	// open and populate hm.db with given fixture
	hm, _ := sql.Open("sqlite3", tempdir+"Database2/hm.db")
	err = hm.Ping()
	if err != nil {
		t.Errorf("unexpected error creating test database: %v", err)
	}

	queryByte, err = os.ReadFile(fixturePath + "hm.sql")
	if err != nil {
		t.Errorf("unexpected error reading from m.db fixture: %v", err)
	}
	query = string(queryByte)

	hm.Exec(query)

	return tempdir
}

// saveStub saves a db.Library struct to a json-formatted stub
func saveStub(t *testing.T, library db.Library, path string) {
	file, err := os.Create(path)
	if err != nil {
		t.Errorf("unexpected error saving library stub: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(library)
	if err != nil {
		t.Errorf("unexpected error saving library stub: %v", err)
	}
}

// loadStub loads a json-formatted db.Library struct stub
func loadStub(t *testing.T, path string) db.Library {
	file, err := os.Open(path)
	if err != nil {
		t.Errorf("unexpected error loading library stub: %v", err)
	}
	defer file.Close()

	var library db.Library
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&library)
	if err != nil {
		t.Errorf("unexpected error loading library stub: %v", err)
	}
	return library
}

// sortSongs sorts db.Library songs based on id,
// used because db.Library.Songs is order-agnostic
func sortSongs(library *db.Library) {
	songs := library.Songs
	sort.Slice(songs, func(i, j int) bool {
		return songs[i].SongID < songs[j].SongID
	})
	library.Songs = songs
}

func TestImportInvalidPath(t *testing.T) {
	_, err := engine.Import("invalid/path", defaultOptions)
	assert.Equal(t, err, errors.New("error initializing database: unable to open database file: no such file or directory"), "Invalid path should throw an error.")
}

func TestImportEmpty(t *testing.T) {
	tempdir := generateDatabase(t, "testdata/fixtures/empty/")
	library, err := engine.Import(tempdir, defaultOptions)
	assert.Nil(t, err, "Empty database import should return no errors.")
	assert.Equal(t, library, db.Library{}, "Empty database import should return an empty library.")
}

func TestImportAlteredPerformanceData(t *testing.T) {
	tempdir := generateDatabase(t, "testdata/fixtures/alteredPerformanceData/")
	library, err := engine.Import(tempdir, defaultOptions)
	sortSongs(&library)
	stub := loadStub(t, "testdata/stubs/alteredPerformanceData.json")
	assert.Nil(t, err, "Valid database import should return no errors.")
	assert.Equal(t, library, stub, "Library should match expected output.")
}

func TestImportSongs(t *testing.T) {
	tempdir := generateDatabase(t, "testdata/fixtures/songs/")
	library, err := engine.Import(tempdir, defaultOptions)
	sortSongs(&library)
	stub := loadStub(t, "testdata/stubs/songs.json")
	assert.Nil(t, err, "Valid database import should return no errors.")
	assert.Equal(t, library, stub, "Library should match expected output.")
}

func TestImportSongsOriginal(t *testing.T) {
	tempdir := generateDatabase(t, "testdata/fixtures/songs/")
	options := defaultOptions
	options.ImportOriginalCues = true
	options.ImportOriginalGrids = true
	library, err := engine.Import(tempdir, options)
	sortSongs(&library)
	stub := loadStub(t, "testdata/stubs/songsOriginal.json")
	assert.Nil(t, err, "Valid database import should return no errors.")
	assert.Equal(t, library, stub, "Library should match expected output.")
}

func TestImportPlaylists(t *testing.T) {
	tempdir := generateDatabase(t, "testdata/fixtures/playlists/")
	library, err := engine.Import(tempdir, defaultOptions)
	sortSongs(&library)
	stub := loadStub(t, "testdata/stubs/playlists.json")
	assert.Nil(t, err, "Valid database import should return no errors.")
	assert.Equal(t, library, stub, "Library should match expected output.")
}

func TestImportNestedPlaylists(t *testing.T) {
	tempdir := generateDatabase(t, "testdata/fixtures/nestedPlaylists/")
	library, err := engine.Import(tempdir, defaultOptions)
	sortSongs(&library)
	stub := loadStub(t, "testdata/stubs/nestedPlaylists.json")
	assert.Nil(t, err, "Valid database import should return no errors.")
	assert.Equal(t, library, stub, "Library should match expected output.")
}

func TestImportCorruptedSong(t *testing.T) {
	tempdir := generateDatabase(t, "testdata/fixtures/corruptSong/")
	library, err := engine.Import(tempdir, defaultOptions)
	sortSongs(&library)
	stub := loadStub(t, "testdata/stubs/corruptSong.json")
	assert.Nil(t, err, "Valid database import should return no errors.")
	assert.Equal(t, library, stub, "Library should match expected output.")
}

func TestImportHistory(t *testing.T) {
	tempdir := generateDatabase(t, "testdata/fixtures/history/")
	library, err := engine.Import(tempdir, defaultOptions)
	sortSongs(&library)
	stub := loadStub(t, "testdata/stubs/history.json")
	assert.Nil(t, err, "Valid database import should return no errors.")
	assert.Equal(t, library, stub, "Library should match expected output.")
}
