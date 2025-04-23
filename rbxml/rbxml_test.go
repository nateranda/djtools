package rbxml_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/nateranda/djtools/lib"
	"github.com/nateranda/djtools/rbxml"
	"github.com/stretchr/testify/assert"
)

var xmlDirExport string = filepath.Join("testdata", "export", "xml")
var jsonDirExport string = filepath.Join("testdata", "export", "json")
var xmlDirImport string = filepath.Join("testdata", "import", "xml")
var jsonDirImport string = filepath.Join("testdata", "import", "json")

type test struct {
	name     string // name of test
	jsonName string // json stub name
	xmlName  string // xml stub name
	saveStub bool   // save a new xml stub or not
}

func loadXml(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("unexpected error reading from XML file at path %s: %v", path, err)
	}

	return data
}

func TestExportInvalidPath(t *testing.T) {
	var library lib.Library
	err := rbxml.Export(&library, "invalid/path.xml")
	assert.Equal(t, errors.New("error exporting library: open invalid/path.xml: no such file or directory"),
		err, "invalid path should throw an error")
}

func TestExport(t *testing.T) {
	tests := []test{
		{"Empty", "empty.json", "empty.xml", false},
		{"Songs", "songs.json", "songs.xml", false},
		{"Playlists", "playlists.json", "playlists.xml", false},
		{"NestedPlaylists", "nestedPlaylists.json", "nestedPlaylists.xml", false},
		{"History", "history.json", "history.xml", false},
		{"CuesLoops", "cuesLoops.json", "cuesLoops.xml", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(jsonDirExport, test.jsonName)
			library := lib.LoadJson(t, path)
			tempPath := filepath.Join(t.TempDir(), "library.xml")
			err := rbxml.Export(&library, tempPath)
			path = filepath.Join(xmlDirExport, test.xmlName)
			if test.saveStub {
				lib.CopyFile(t, tempPath, path)
				t.Fail()
			}
			export := loadXml(t, tempPath)
			check := loadXml(t, path)
			assert.Nil(t, err, "Valid database import should return no errors.")
			assert.Equal(t, check, export, "Library should match expected output.")
		})
	}
}

func TestImportInvalidPath(t *testing.T) {
	_, err := rbxml.Import("invalid/path/library.xml")
	assert.Equal(t, errors.New("error reading file: open invalid/path/library.xml: no such file or directory"),
		err, "Invalid path should throw an error.")
}

func TestImport(t *testing.T) {
	tests := []test{
		{"Empty", "empty.json", "empty.xml", false},
		{"Songs", "songs.json", "songs.xml", false},
		{"CorruptSong", "corruptSong.json", "corruptSong.xml", false},
		{"CuesLoops", "cuesLoops.json", "cuesLoops.xml", false},
		{"Playlists", "playlists.json", "playlists.xml", false},
		{"NestedPlaylists", "nestedPlaylists.json", "nestedPlaylists.xml", false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(xmlDirImport, test.xmlName)
			library, err := rbxml.Import(path)
			library.SortSongs()
			path = filepath.Join(jsonDirImport, test.jsonName)
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
