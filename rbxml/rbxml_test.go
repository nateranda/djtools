package rbxml_test

import (
	"errors"
	"os"
	"testing"

	"github.com/nateranda/djtools/lib"
	"github.com/nateranda/djtools/rbxml"
	"github.com/stretchr/testify/assert"
)

const (
	xmlDirExport  string = "testdata/export/xml/"
	jsonDirExport string = "testdata/export/json/"
	xmlDirImport  string = "testdata/import/xml/"
	jsonDirImport string = "testdata/import/json/"
)

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
			library := lib.LoadJson(t, jsonDirExport+test.jsonName)
			tempdir := t.TempDir() + "/"
			err := rbxml.Export(&library, tempdir+"library.xml")
			if test.saveStub {
				lib.CopyFile(t, tempdir+"library.xml", xmlDirExport+test.xmlName)
				t.Fail()
			}
			export := loadXml(t, tempdir+"library.xml")
			check := loadXml(t, xmlDirExport+test.xmlName)
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
			library, err := rbxml.Import(xmlDirImport + test.xmlName)
			library.SortSongs()
			if test.saveStub {
				lib.SaveJson(t, library, jsonDirImport+test.jsonName)
				t.Fail()
			}
			stub := lib.LoadJson(t, jsonDirImport+test.jsonName)
			assert.Nil(t, err, "Valid database import should return no errors.")
			assert.Equal(t, library, stub, "Library should match expected output.")
		})
	}
}
