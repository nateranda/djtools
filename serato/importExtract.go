package serato

import (
	"encoding/binary"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf16"

	"github.com/dhowden/tag"
)

func importExtract(path string) error {
	crates, err := importExtractCrates(path)
	if err != nil {
		return err
	}
	songs, err := importExtractGeobs(crates)
	if err != nil {
		return err
	}
	fmt.Printf("songs: %v\n", songs)
	return nil
}

func importExtractCrates(path string) ([]crate, error) {
	var crates []crate
	paths, err := listCrateFiles(path)
	if err != nil {
		return nil, err
	}
	for _, path := range paths {
		crate, err := importExtractCrate(path)
		if err != nil {
			return nil, err
		}
		crates = append(crates, crate)
	}
	return crates, nil
}

func importExtractCrate(path string) (crate, error) {
	var entries crate
	var err error

	pathWithoutExtension := path[:len(path)-len(filepath.Ext(path))]
	entries.filename = filepath.Base(pathWithoutExtension)

	file, err := os.ReadFile(path)
	if err != nil {
		return crate{}, err
	}

	for len(file) > 0 {
		file, err = entries.extractCrateEntry(file)
		if err != nil {
			return crate{}, err
		}
	}

	return entries, nil
}

func importExtractGeobs(crates []crate) ([]song, error) {
	var songs []song
	for _, crate := range crates {
		for _, path := range crate.paths {
			var song song
			song.path = path
			geobs, err := importExtractGeob(path)
			if err != nil {
				return nil, err
			}
			song.geobs = geobs
			songs = append(songs, song)
		}
	}
	return songs, nil
}

func importExtractGeob(path string) ([]geob, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	metadata, err := tag.ReadFrom(file)
	if err != nil {
		return nil, fmt.Errorf("error reading metadata: %v", err)
	}

	raw := metadata.Raw()
	if raw == nil {
		return nil, fmt.Errorf("no raw metadata found")
	}
	var geobs []geob
	for key := range raw {
		if strings.HasPrefix(key, "GEOB") {
			if value, ok := raw[key].([]byte); ok {
				var geob geob
				geob.name = key
				geob.value = value
				geobs = append(geobs, geob)
			}
		}
	}

	sort.Slice(geobs, func(i, j int) bool {
		return geobs[i].name < geobs[j].name
	})

	return geobs, nil
}

func (c *crate) extractCrateEntry(file []byte) ([]byte, error) {
	// extract key (4 bytes)
	if len(file) < 4 {
		return nil, fmt.Errorf("file too short to extract key")
	}
	key := string(file[:4])
	file = file[4:]

	// extract length (4 bytes)
	length := int(binary.BigEndian.Uint32(file[:+4]))
	file = file[4:]

	// extract value based on key type (variable bytes)
	if len(file) < length {
		return nil, fmt.Errorf("file too short to extract value")
	}
	valueBytes := file[:length]
	file = file[length:]

	switch key {
	case "otrk":
		var err error
		valueBytes, err = c.extractCrateEntry(valueBytes)
		if len(valueBytes) != 0 {
			return nil, fmt.Errorf("invalid otrk entry")
		}
		if err != nil {
			return nil, err
		}
		return file, nil
	case "vrsn":
		version, err := utf16ToString(valueBytes)
		if err != nil {
			return nil, err
		}
		c.version = version
		return file, nil
	case "ptrk":
		path, err := utf16ToString(valueBytes)
		if err != nil {
			return nil, err
		}
		c.paths = append(c.paths, string(filepath.Separator)+path)
		return file, nil
	}

	// skip unneeded entries
	switch key[:1] {
	case "o", "t", "p", "u", "s", "b":
		return file, nil
	}

	return nil, fmt.Errorf("key %s not of supported type", key)
}

func utf16ToString(data []byte) (string, error) {
	// check for even input
	if len(data)%2 != 0 {
		return "", fmt.Errorf("invalid UTF-16 byte slice length: %d", len(data))
	}

	uint16s := make([]uint16, len(data)/2)
	for i := 0; i < len(data); i += 2 {
		uint16s[i/2] = binary.BigEndian.Uint16(data[i : i+2])
	}

	runes := utf16.Decode(uint16s)
	content := string(runes)
	return content, nil
}

func listCrateFiles(path string) ([]string, error) {
	var files []string
	cratepath := filepath.Join(path, "Subcrates")
	entries, err := os.ReadDir(cratepath)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && filepath.Ext(entry.Name()) == ".crate" {
			files = append(files, filepath.Join(cratepath, entry.Name()))
		}
	}

	return files, nil
}
