package db

import (
	"fmt"
)

type Marker struct {
	StartPosition float64
	Bpm           float64
	BeatNumber    int
}

type HotCue struct {
	Name     string
	Offset   float64
	Position int
	Color    string
}

type Loop struct {
	Name     string
	Start    float64
	End      float64
	Position int
	Color    string
}

type Song struct {
	SongID       int
	Title        string
	Artist       string
	Composer     string
	Album        string
	Grouping     string
	Genre        string
	Filetype     string
	Size         int
	Length       float32
	TrackNumber  int
	Year         int
	Bpm          float32
	DateModified int
	DateAdded    int
	Bitrate      int
	SampleRate   float64
	Comment      string
	PlayCount    int
	LastPlayed   int
	Rating       int
	Path         string
	Remixer      string
	Key          string
	Label        string
	Mix          string
	Color        string
	Cue          float64
	Grid         []Marker
	Cues         []HotCue
	Loops        []Loop
}

type Playlist struct {
	PlaylistID   int
	Position     int
	Name         string
	Songs        []*Song
	SubPlaylists []Playlist
}

type Library struct {
	Songs     []Song
	Playlists []Playlist
}

// getSong takes a Song slice and returns a pointer to the Song with the given id
func GetSong(songList []Song, id int) (*Song, error) {
	for i := range songList {
		if songList[i].SongID == id {
			return &songList[i], nil
		}
	}
	return nil, fmt.Errorf("NotFoundError: did not find a Song matching %d", id)
}
