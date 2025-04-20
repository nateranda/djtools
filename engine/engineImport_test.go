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

const (
	fixturesDir = "testdata/fixtures/"
	stubsDir    = "testdata/stubs/"
)

var defaultOptions = engine.ImportOptions{
	PreserveOriginalPaths: true, // preserve original relative filepaths for operating system parity
}

type test struct {
	name     string               // name of test
	dirname  string               // fixture directory name
	filename string               // stub file name
	saveStub bool                 // save a new stub or not
	options  engine.ImportOptions // importOptions to pass
}

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
	assert.Equal(t, err,
		errors.New("error initializing database: unable to open database file: no such file or directory"),
		"Invalid path should throw an error.")
}

func TestImport(t *testing.T) {
	tests := []test{
		{"Empty", "empty/", "empty.json", false, defaultOptions},
		{"Songs", "songs/", "songs.json", false, defaultOptions},
		{"SongsOriginal", "songsOriginal/", "songsOriginal.json", false, engine.ImportOptions{
			PreserveOriginalPaths: true,
			ImportOriginalCues:    true,
			ImportOriginalGrids:   true,
		}},
		{"AlteredPerformanceData", "alteredPerformanceData/", "alteredPerformanceData.json", false, defaultOptions},
		{"Playlists", "playlists/", "playlists.json", false, defaultOptions},
		{"NestedPlaylists", "nestedPlaylists/", "nestedPlaylists.json", false, defaultOptions},
		{"CorruptSong", "corruptSong/", "corruptSong.json", false, defaultOptions},
		{"History", "history/", "history.json", false, defaultOptions},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tempdir := generateDatabase(t, fixturesDir+test.dirname)
			library, err := engine.Import(tempdir, test.options)
			if test.saveStub {
				saveStub(t, library, stubsDir+test.filename)
			}
			sortSongs(&library)
			stub := loadStub(t, stubsDir+test.filename)
			assert.Nil(t, err, "Valid database import should return no errors.")
			assert.Equal(t, library, stub, "Library should match expected output.")
		})
	}
}
