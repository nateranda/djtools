package lib

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
)

// Marker is a marker in a beatgrid.
type Marker struct {
	StartPosition float64 // start position in seconds
	Bpm           float64 // beats per minute
	BeatNumber    int     // beat in measure to start on, 0-indexed, assumed 4/4
}

// HotCue is a saved hot cue for a song.
type HotCue struct {
	Name     string  // name of hot cue
	Offset   float64 // offset of hot cue in seconds
	Position int     // position of hot cue, 1-indexed
	Color    string  // color of hot cue button, hex code
}

// Loop is a saved loop for a song.
type Loop struct {
	Name     string  // name of loop
	Start    float64 // start time of loop, seconds
	End      float64 // end time of loop, seconds
	Position int     // position of loop, 1-indexed
	Color    string  // color of loop button, hex code
}

// Song is the metadata, analysis data, and saved cues/loops for a song.
type Song struct {
	SongID       int      // song id used by software
	Title        string   // title
	Artist       string   // artist
	Composer     string   // composer
	Album        string   // album song is from
	Grouping     string   // grouping
	Genre        string   // genre
	Filetype     string   // filetype, abbreviated lowercase
	Size         int      // file size, bytes
	Length       float32  // song length, seconds
	TrackNumber  int      // number in album
	Year         int      // release year
	Bpm          float32  // beats per minute
	DateModified int      // date last modified, unix
	DateAdded    int      // date added to library, unix
	Bitrate      int      // bitrate, kbps
	SampleRate   float64  // sample rate, hz
	Comment      string   // comment
	PlayCount    int      // play count
	LastPlayed   int      // date last played, unix
	Rating       int      // rating in multiples of 20: 0*=0, 1*=20... 5*=100
	Path         string   // song absolute path
	Remixer      string   // remixer
	Key          int      // key in int representation of camelot, 0-indexed: 0=8B, 1=8A, 2=9B... 23=7A
	Label        string   // label
	Mix          string   // mix
	Color        string   // color, hex code
	Cue          float64  // cue location, seconds
	Grid         []Marker // slice of Marker structs, ordered by start position
	Cues         []HotCue // slice of Cue structs, unordered
	Loops        []Loop   // slice of Loop structs, unordered
	Corrupt      bool     // is the song file corrupted?
}

// Playlist is a set of ordered songs which can contain other playlists.
// A folder is just a Playlist with no songs that contains other playlists.
type Playlist struct {
	PlaylistID   int        // playlist id used by software
	Name         string     // name of playlist
	Songs        []int      // slice of song ids in order
	SubPlaylists []Playlist // slice of child playlists in order, can be recursive
}

// Library is the entire library of a DJ software.
type Library struct {
	Songs     []Song     // slice of Song structs, unordered
	Playlists []Playlist // slice of Playlist structs, ordered by position
}

// Save saves a Library struct to a json file.
// Used for development and testing purposes only.
func (library *Library) Save(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("unexpected error saving library stub: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(library)
	if err != nil {
		return fmt.Errorf("unexpected error saving library stub: %v", err)
	}
	return nil
}

// Load loads a Library struct from a json file.
// Used for development and testing purposes only.
func (library *Library) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unexpected error loading library stub: %v", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(library)
	if err != nil {
		return fmt.Errorf("unexpected error loading library stub: %v", err)
	}
	return nil
}

// SortSongs sorts a Library's Songs based on id and sorts each Song's Cues
// and Loops by position, used because Library.Songs is order-agnostic
func (l *Library) SortSongs() {
	// sort songs by id
	sort.Slice(l.Songs, func(i, j int) bool {
		return l.Songs[i].SongID < l.Songs[j].SongID
	})

	// sort cues and loops by position
	for i := range l.Songs {
		sort.Slice(l.Songs[i].Cues, func(a, b int) bool {
			return l.Songs[i].Cues[a].Position < l.Songs[i].Cues[b].Position
		})
		sort.Slice(l.Songs[i].Loops, func(a, b int) bool {
			return l.Songs[i].Loops[a].Position < l.Songs[i].Loops[b].Position
		})
	}
}

// CheckCorruptedSongs removes songs marked as corrupted from the Library
func (l *Library) CheckCorruptedSongs() {
	for i, song := range l.Songs {
		// this is expensive, but it should happen rarely so it's ok
		if song.Corrupt {
			// remove song from library.Songs (doesn't preserve order)
			l.Songs[i] = l.Songs[len(l.Songs)-1]
			l.Songs = l.Songs[:len(l.Songs)-1]

			// remove song from playlists
			l.Playlists = removeSongFromPlaylists(l.Playlists, song.SongID)
		}
	}
}

func removeSongFromPlaylists(playlists []Playlist, songID int) []Playlist {
	for i := range playlists {
		var updatedSongIDs []int
		for _, id := range playlists[i].Songs {
			if id != songID {
				updatedSongIDs = append(updatedSongIDs, id)
			}
		}
		playlists[i].Songs = updatedSongIDs

		if playlists[i].SubPlaylists != nil {
			playlists[i].SubPlaylists = removeSongFromPlaylists(playlists[i].SubPlaylists, songID)
		}
	}

	return playlists
}
