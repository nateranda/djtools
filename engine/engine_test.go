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

var defaultOptions = engine.ImportOptions{PreserveOriginalPaths: true}

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

func saveStub(t *testing.T, library db.Library, path string) {
	file, err := os.Create(path)
	if err != nil {
		t.Errorf("unexpected error saving library stub: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(library)
	if err != nil {
		t.Errorf("unexpected error saving library stub: %v", err)
	}
}

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

// db.Library.Songs is order-agnostic, so songs will be sorted by id for testing purposes
func sortSongs(library *db.Library) {
	songs := library.Songs
	sort.Slice(songs, func(i, j int) bool {
		return songs[i].SongID < songs[j].SongID
	})
	library.Songs = songs
}

func TestImportEmpty(t *testing.T) {
	tempdir := generateDatabase(t, "testdata/fixtures/empty/")
	library, err := engine.Import(tempdir, defaultOptions)
	assert.Nil(t, err, "Empty database import should return no errors.")
	assert.Equal(t, library, db.Library{}, "Empty database import should return an empty library.")
}
func TestImportInvalidPath(t *testing.T) {
	_, err := engine.Import("invalid/path", defaultOptions)
	assert.Equal(t, err, errors.New("error initializing database: unable to open database file: no such file or directory"), "Invalid path should throw an error.")
}

func TestImportAlteredPerformanceData(t *testing.T) {
	tempdir := generateDatabase(t, "testdata/fixtures/alteredPerformanceData/")
	library, err := engine.Import(tempdir, defaultOptions)
	sortSongs(&library)
	stub := loadStub(t, "testdata/stubs/alteredPerformanceData.json")
	assert.Nil(t, err, "Valid database import should return no errors.")
	assert.Equal(t, library, stub, "Library should match expected output.")
}
