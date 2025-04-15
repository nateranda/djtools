package rbxml

import (
	"fmt"
	"net/url"
	"path/filepath"
	"time"

	"github.com/nateranda/djtools/db"
)

func unixToDate(date int) string {
	t := time.Unix(int64(date), 0)
	return t.Format("2006-01-02")
}

func pathToURI(path string) (string, error) {
	// Convert the path to an absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("error converting path to URI: %v", err)
	}
	uriPath := filepath.ToSlash(absPath)
	return "file://localhost/" + url.PathEscape(uriPath), nil
}

func exportConvertRating(rating int) (int32, error) {
	switch rating {
	case 0:
		return 0, nil
	case 20:
		return 51, nil
	case 40:
		return 102, nil
	case 60:
		return 153, nil
	case 80:
		return 204, nil
	case 100:
		return 255, nil
	}
	return 0, fmt.Errorf("NoMatchError: rating %d did not match convention. Must be 0, 20, 40, 60, 80, or 100", rating)
}

func exportConvertSong(library *db.Library) ([]track, error) {
	var tracks []track
	for _, song := range library.Songs {
		rating, err := exportConvertRating(song.Rating)
		if err != nil {
			return nil, fmt.Errorf("error converting song: %v", err)
		}
		path, err := pathToURI(song.Path)
		if err != nil {
			return nil, fmt.Errorf("error converting song: %v", err)
		}
		tracks = append(tracks, track{
			trackId:      song.SongID,
			name:         song.Title,
			artist:       song.Artist,
			composer:     song.Composer,
			album:        song.Album,
			grouping:     song.Grouping,
			genre:        song.Genre,
			kind:         song.Filetype,
			size:         int64(song.Size),
			totalTime:    float64(song.Length), // make sure this is rounded?
			discNumber:   0,
			trackNumber:  int32(song.TrackNumber),
			year:         int32(song.Year),
			averageBpm:   float64(song.Bpm),
			dateModified: unixToDate(song.DateModified),
			dateAdded:    unixToDate(song.DateAdded),
			bitRate:      int32(song.Bitrate),
			sampleRate:   song.SampleRate,
			comments:     song.Comment,
			playCount:    int32(song.PlayCount),
			lastPlayed:   unixToDate(song.LastPlayed),
			rating:       rating,
			location:     path,
			remixer:      song.Remixer,
			tonality:     song.Key, // wrong key type?
			label:        song.Label,
			mix:          song.Mix,
			colour:       song.Color,
		})
	}
	return tracks, nil
}

func exportConvert(library *db.Library) (collection, error) {
	var collection collection
	var err error

	collection.tracks, err = exportConvertSong(library)
	if err != nil {
		return collection, err
	}
	return collection, nil
}
