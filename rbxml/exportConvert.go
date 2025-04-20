package rbxml

import (
	"fmt"
	"net/url"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/nateranda/djtools/db"
)

func unixToDate(date int) string {
	t := time.Unix(int64(date), 0)
	return t.Format("2006-01-02")
}

func pathToURI(path string) string {
	// this is a jank fix, replace with something more robust?
	uriPath := filepath.ToSlash(path)
	uriPath = url.PathEscape(uriPath)
	uriPath = strings.ReplaceAll(uriPath, "%2F", "/") // keep slashes
	return "file://localhost/" + uriPath
}

func exportConvertTonality(key int) (string, error) {
	switch key {
	case 0:
		return "8B", nil
	case 1:
		return "8A", nil
	case 2:
		return "9B", nil
	case 3:
		return "9A", nil
	case 4:
		return "10B", nil
	case 5:
		return "10A", nil
	case 6:
		return "11B", nil
	case 7:
		return "11A", nil
	case 8:
		return "12B", nil
	case 9:
		return "12A", nil
	case 10:
		return "1B", nil
	case 11:
		return "1A", nil
	case 12:
		return "2B", nil
	case 13:
		return "2A", nil
	case 14:
		return "3B", nil
	case 15:
		return "3A", nil
	case 16:
		return "4B", nil
	case 17:
		return "4A", nil
	case 18:
		return "5B", nil
	case 19:
		return "5A", nil
	case 20:
		return "6B", nil
	case 21:
		return "6A", nil
	}
	return "", fmt.Errorf("key outside accepted range")
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

func exportConvertPositionMarks(song *db.Song) []positionMark {
	var positionMarks []positionMark
	// add cue point
	positionMarks = append(positionMarks, positionMark{
		MarkType: 0,
		Start:    song.Cue,
		Num:      -1,
	})

	// add hot cues
	for _, cue := range song.Cues {
		positionMarks = append(positionMarks, positionMark{
			Name:     cue.Name,
			MarkType: 0,
			Start:    cue.Offset,
			Num:      int32(cue.Position - 1),
		})
	}

	// add loops
	for _, loop := range song.Loops {
		positionMarks = append(positionMarks, positionMark{
			Name:     loop.Name,
			MarkType: 4,
			Start:    loop.Start,
			End:      loop.End,
			Num:      int32(loop.Position - 1),
		})
	}

	return positionMarks
}

func exportConvertGrid(song *db.Song) []tempo {
	var tempos []tempo

	for _, grid := range song.Grid {
		tempos = append(tempos, tempo{
			Inizio:  grid.StartPosition,
			Bpm:     grid.Bpm,
			Metro:   "4/4", // assumed, may add time signature support later
			Battito: int32(grid.BeatNumber) + 1,
		})
	}

	return tempos
}

func exportConvertSubPlaylists(playlist db.Playlist) []node {
	var nodes []node

	// add playlist node containing tracks
	if playlist.Songs != nil {
		var tracks []nodeTrack
		for id := range playlist.Songs {
			tracks = append(tracks, nodeTrack{Id: int32(id)})
		}
		nodes = append(nodes, node{
			NodeType: 1,
			Name:     playlist.Name,
			Entries:  int32(len(playlist.Songs)),
			Tracks:   &tracks,
		})
	}

	// add folder node containing sub-playlists
	if playlist.SubPlaylists != nil {
		var subNodes []node
		// add sub-playlist nodes recursively
		for _, playlist := range playlist.SubPlaylists {
			subNodes = slices.Concat(subNodes, exportConvertSubPlaylists(playlist))
		}
		nodes = append(nodes, node{
			NodeType: 0,
			Name:     playlist.Name,
			Count:    int32(len(playlist.SubPlaylists)),
			Nodes:    &subNodes,
		})
	}

	return nodes
}

func exportConvertPlaylist(library *db.Library) node {
	var nodes []node
	for _, playlist := range library.Playlists {
		nodes = slices.Concat(nodes, exportConvertSubPlaylists(playlist))
	}

	return node{
		NodeType: 0,
		Name:     "ROOT",
		Count:    int32(len(library.Playlists)),
		Nodes:    &nodes,
	}
}

func exportConvertSong(library *db.Library) ([]track, error) {
	var tracks []track
	for _, song := range library.Songs {
		rating, err := exportConvertRating(song.Rating)
		if err != nil {
			return nil, fmt.Errorf("error converting song: %v", err)
		}
		path := pathToURI(song.Path)
		tonality, err := exportConvertTonality(song.Key)
		if err != nil {
			return nil, fmt.Errorf("error converting song: %v", err)
		}
		positionMarks := exportConvertPositionMarks(&song)
		tempos := exportConvertGrid(&song)
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
			Tonality:     tonality, // wrong key type?
			Label:        song.Label,
			Mix:          song.Mix,
			Colour:       song.Color,
			PositionMark: &positionMarks,
			Tempo:        &tempos,
		})
	}
	return tracks, nil
}

func exportConvert(library *db.Library) (djPlaylists, error) {
	djPlaylists := djPlaylists{
		Version: "1.0.0",
		Product: product{
			Name:    "djtools",
			Version: version,
			Company: "djtools",
		},
	}
	var err error

	djPlaylists.Collection.Tracks, err = exportConvertSong(library)
	djPlaylists.Collection.Entries = int32(len(djPlaylists.Collection.Tracks))
	if err != nil {
		return djPlaylists, err
	}

	djPlaylists.Playlists = playlists{Node: exportConvertPlaylist(library)}

	return djPlaylists, nil
}
