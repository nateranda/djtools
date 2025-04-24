package serato

import (
	"fmt"
	"os"
	"sort"
	"strings"

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
