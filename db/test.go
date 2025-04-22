package db

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"testing"
)

// SaveJson saves a db.Library struct to a json-formatted stub
func SaveJson(t *testing.T, library Library, path string) {
	t.Helper()
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

// LoadJson loads a json-formatted db.Library struct stub
func LoadJson(t *testing.T, path string) Library {
	t.Helper()
	file, err := os.Open(path)
	if err != nil {
		t.Errorf("unexpected error loading library stub: %v", err)
	}
	defer file.Close()

	var library Library
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&library)
	if err != nil {
		t.Errorf("unexpected error loading library stub: %v", err)
	}

	return library
}

// CopyFile copies a file from a source to a destination
func CopyFile(t *testing.T, srcPath string, destPath string) {
	t.Helper()
	srcFile, err := os.Open(srcPath)
	if err != nil {
		t.Errorf("unexpected error opening source file: %v", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		t.Errorf("unexpected error creating destination file: %v", err)
	}

	i, err := io.Copy(destFile, srcFile)
	if err != nil {
		t.Errorf("unexpected error copying source file: %v", err)
	}
	fmt.Println(i)

	err = destFile.Sync()
	if err != nil {
		t.Errorf("unexpected error syncing destination file: %v", err)
	}
}
