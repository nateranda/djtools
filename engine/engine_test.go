package engine_test

import (
	"database/sql"
	"errors"
	"os"
	"sort"
	"testing"

	"github.com/nateranda/djtools/engine"
	"github.com/nateranda/djtools/lib"
	"github.com/stretchr/testify/assert"

	_ "github.com/mattn/go-sqlite3"
)

const (
	fixturesDir = "testdata/import/fixtures/"
	stubsDir    = "testdata/import/stubs/"
)

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

// sortSongs sorts lib.Library songs based on id and sorts each song's cues
// and loops by position, used because lib.Library.Songs is order-agnostic
func sortSongs(library *lib.Library) {
	songs := library.Songs

	// sort songs by id
	sort.Slice(songs, func(i, j int) bool {
		return songs[i].SongID < songs[j].SongID
	})

	// sort cues and loops by position
	for i := range songs {
		sort.Slice(songs[i].Cues, func(a, b int) bool {
			return songs[i].Cues[a].Position < songs[i].Cues[b].Position
		})
		sort.Slice(songs[i].Loops, func(a, b int) bool {
			return songs[i].Loops[a].Position < songs[i].Loops[b].Position
		})
	}

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
			sortSongs(&library)
			if test.saveStub {
				lib.SaveJson(t, library, stubsDir+test.filename)
				t.Fail()
			}
			stub := lib.LoadJson(t, stubsDir+test.filename)
			assert.Nil(t, err, "Valid database import should return no errors.")
			assert.Equal(t, library, stub, "Library should match expected output.")
		})
	}
}
