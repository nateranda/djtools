package serato

import (
	"encoding/binary"
	"fmt"
	"os"
	"unicode/utf16"
)

type crate struct {
	version string
	paths   []string
}

func ExtractCrate(path string) (crate, error) {
	var entries crate
	var err error

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
		c.paths = append(c.paths, path)
		return file, nil
	}

	// return empty entry for unneeded key
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
