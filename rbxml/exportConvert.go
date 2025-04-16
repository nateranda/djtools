package rbxml

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/nateranda/djtools/db"
)

func unixToDate(date int) string {
	t := time.Unix(int64(date), 0)
	return t.Format("2006-01-02")
}

func pathToURI(path string) (string, error) {
	// this is a jank fix, replace with something more robust?
	uriPath := filepath.ToSlash(path)
	uriPath = url.PathEscape(uriPath)
	uriPath = strings.ReplaceAll(uriPath, "%2F", "/") // keep slashes
	return "file://localhost/" + uriPath, nil
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
			TrackId:      song.SongID,
			Name:         song.Title,
			Artist:       song.Artist,
			Composer:     song.Composer,
			Album:        song.Album,
			Grouping:     song.Grouping,
			Genre:        song.Genre,
			Kind:         song.Filetype,
			Size:         int64(song.Size),
			TotalTime:    float64(song.Length), // make sure this is rounded?
			DiscNumber:   0,
			TrackNumber:  int32(song.TrackNumber),
			Year:         int32(song.Year),
			AverageBpm:   float64(song.Bpm),
			DateModified: unixToDate(song.DateModified),
			DateAdded:    unixToDate(song.DateAdded),
			BitRate:      int32(song.Bitrate),
			SampleRate:   song.SampleRate,
			Comments:     song.Comment,
			PlayCount:    int32(song.PlayCount),
			LastPlayed:   unixToDate(song.LastPlayed),
			Rating:       rating,
			Location:     path,
			Remixer:      song.Remixer,
			Tonality:     song.Key, // wrong key type?
			Label:        song.Label,
			Mix:          song.Mix,
			Colour:       song.Color,
		})
	}
	return tracks, nil
}

func exportConvert(library *db.Library) (djPlaylists, error) {
	var djPlaylists djPlaylists
	var err error

	djPlaylists.Collection.Tracks, err = exportConvertSong(library)
	djPlaylists.Collection.Entries = int32(len(djPlaylists.Collection.Tracks))
	if err != nil {
		return djPlaylists, err
	}
	return djPlaylists, nil
}
