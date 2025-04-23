package engine_test

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/nateranda/djtools/engine"
	"github.com/nateranda/djtools/lib"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

var fixturesDir string = filepath.Join("testdata", "import", "fixtures")
var stubsDir string = filepath.Join("testdata", "import", "stubs")

var defaultOptions = engine.ImportOptions{
	// preserve original relative filepaths for operating system parity
	PreserveOriginalPaths: true,
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
	t.Helper()
	tempdir := t.TempDir()

	//make Database2 directory inside of temp directory
	path := filepath.Join(tempdir, "Database2")
	os.Mkdir(path, 0755)

	// open and populate m.db with given fixture
	path = filepath.Join(tempdir, "Database2", "m.db")
	m, _ := sql.Open("sqlite3", path)
	err := m.Ping()
	if err != nil {
		t.Errorf("unexpected error creating test database: %v", err)
	}

	path = filepath.Join(fixturePath, "m.sql")
	queryByte, err := os.ReadFile(path)
	if err != nil {
		t.Errorf("unexpected error reading from m.db fixture: %v", err)
	}
	query := string(queryByte)

	m.Exec(query)

	// open and populate hm.db with given fixture
	path = filepath.Join(tempdir, "Database2", "hm.db")
	hm, _ := sql.Open("sqlite3", path)
	err = hm.Ping()
	if err != nil {
		t.Errorf("unexpected error creating test database: %v", err)
	}

	path = filepath.Join(fixturePath, "hm.sql")
	queryByte, err = os.ReadFile(path)
	if err != nil {
		t.Errorf("unexpected error reading from m.db fixture: %v", err)
	}
	query = string(queryByte)

	hm.Exec(query)

	return tempdir
}

func TestImportInvalidPath(t *testing.T) {
	_, err := engine.Import("invalid/path", defaultOptions)
	assert.Equal(t, err,
		errors.New("error initializing database: unable to open database file: no such file or directory"),
		"Invalid path should throw an error.")
}

func TestImport(t *testing.T) {
	tests := []test{
		{"Empty", "empty", "empty.json", false, defaultOptions},
		{"Songs", "songs", "songs.json", false, defaultOptions},
		{"SongsOriginal", "songsOriginal", "songsOriginal.json", false, engine.ImportOptions{
			PreserveOriginalPaths: true,
			ImportOriginalCues:    true,
			ImportOriginalGrids:   true,
		}},
		{"AlteredPerformanceData", "alteredPerformanceData", "alteredPerformanceData.json", false, defaultOptions},
		{"Playlists", "playlists", "playlists.json", false, defaultOptions},
		{"NestedPlaylists", "nestedPlaylists", "nestedPlaylists.json", false, defaultOptions},
		{"CorruptSong", "corruptSong", "corruptSong.json", false, defaultOptions},
		{"History", "history", "history.json", false, defaultOptions},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(fixturesDir, test.dirname)
			tempdir := generateDatabase(t, path)
			library, err := engine.Import(tempdir, test.options)
			library.SortSongs()
			path = filepath.Join(stubsDir, test.filename)
			if test.saveStub {
				lib.SaveJson(t, library, path)
				t.Fail()
			}
			stub := lib.LoadJson(t, path)
			assert.Nil(t, err, "Valid database import should return no errors.")
			assert.Equal(t, library, stub, "Library should match expected output.")
		})
	}
}
