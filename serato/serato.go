package serato

import (
	"encoding/binary"
	"fmt"
	"os"
	"unicode/utf16"
)

type crateEntry struct {
	key   string
	value any
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

func extractCrateEntries(file []byte) ([]crateEntry, error) {
	var crate []crateEntry
	var err error
	for len(file) > 0 {
		var entry crateEntry
		entry, file, err = extractCrateEntry(file)
		if err != nil {
			return nil, err
		}
		crate = append(crate, entry)
	}
	return crate, nil
}

func extractCrateEntry(file []byte) (crateEntry, []byte, error) {
	var crate crateEntry

	// extract key (4 bytes)
	if len(file) < 4 {
		return crateEntry{}, nil, fmt.Errorf("file too short to extract key")
	}
	crate.key = string(file[:4])
	file = file[4:]

	// extract length (4 bytes)
	length := int(binary.BigEndian.Uint32(file[:+4]))
	file = file[4:]

	// extract value based on key type (variable bytes)
	if len(file) < length {
		return crateEntry{}, nil, fmt.Errorf("file too short to extract value")
	}
	valueBytes := file[:length]
	file = file[length:]

	switch crate.key {
	case "otrk":
		var err error
		crate.value, err = extractCrateEntries(valueBytes)
		if err != nil {
			return crateEntry{}, nil, err
		}
		return crate, file, nil
	case "vrsn", "ptrk":
		value, err := utf16ToString(valueBytes)
		if err != nil {
			return crateEntry{}, nil, err
		}
		crate.value = value
		return crate, file, nil
	}

	// return empty entry for unneeded key
	switch crate.key[:1] {
	case "o", "t", "p", "u", "s", "b":
		return crateEntry{}, file, nil
	}

	return crateEntry{}, nil, fmt.Errorf("key %s not of supported type", crate.key)
}

func ExtractCrate(path string) ([]crateEntry, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	entries, err := extractCrateEntries(file)
	if err != nil {
		return nil, err
	}
	return entries, nil
}
