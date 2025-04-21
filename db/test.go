package db

import (
	"encoding/json"
	"os"
	"testing"
)

// saveStub saves a db.Library struct to a json-formatted stub
func SaveStub(t *testing.T, library Library, path string) {
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
func LoadStub(t *testing.T, path string) Library {
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
