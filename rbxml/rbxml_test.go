package rbxml_test

import (
	"os"
	"testing"

	"github.com/nateranda/djtools/db"
	"github.com/nateranda/djtools/rbxml"
	"github.com/stretchr/testify/assert"
)

const (
	xmlDirExport  string = "testdata/export/xml/"
	jsonDirExport string = "testdata/export/json/"
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

func TestExport(t *testing.T) {
	library := db.LoadJson(t, jsonDirExport+"cuesLoops.json")
	tempdir := t.TempDir() + "/"
	err := rbxml.Export(&library, tempdir+"library.xml")
	db.CopyFile(t, tempdir+"library.xml", xmlDirExport+"cuesLoops.xml")
	export := loadXml(t, tempdir+"library.xml")
	check := loadXml(t, xmlDirExport+"cuesLoops.xml")
	assert.Nil(t, err, "Valid database import should return no errors.")
	assert.Equal(t, check, export, "Library should match expected output.")
}
