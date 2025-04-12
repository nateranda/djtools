package db

import (
	"fmt"
)

// unused
type ImportOptions struct {
}

// unused
type ExportOptions struct {
}

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
	PlaylistID int
	Position   int
	Name       string
	Songs      []*Song
}

type Folder struct {
	FolderID  int
	Name      string
	Position  int
	Playlists []*Playlist
}

type Library struct {
	Songs     []Song
	Playlists []Playlist
	Folders   []Folder
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

// getPlaylist takes a Playlist slice and returns a pointer to the Playlist with the given id
func GetPlaylist(playlistList []Playlist, id int) (*Playlist, error) {
	for i := range playlistList {
		if playlistList[i].PlaylistID == id {
			return &playlistList[i], nil
		}
	}
	return nil, fmt.Errorf("NotFoundError: did not find a Playlist matching %d", id)
}

// getFolder takes a Folder slice and returns a pointer to the Folder with the given id
func GetFolder(folderList []Folder, id int) (*Folder, error) {
	for i := range folderList {
		if folderList[i].FolderID == id {
			return &folderList[i], nil
		}
	}
	return nil, fmt.Errorf("NotFoundError: did not find a Folder matching %d", id)
}

func TestLibrary() Library {
	songs := []Song{
		{SongID: 1, Title: "test1", Artist: "artist1"},
		{SongID: 2, Title: "test2", Artist: "artist2"},
		{SongID: 3, Title: "test3", Artist: "artist3"},
	}
	playlists := []Playlist{
		{PlaylistID: 1, Position: 1, Name: "playlist1", Songs: []*Song{&songs[0], &songs[1], &songs[2]}},
		{PlaylistID: 2, Position: 2, Name: "playlist2", Songs: []*Song{&songs[0], &songs[1]}},
	}
	folders := []Folder{{FolderID: 1, Name: "folder1", Position: 1, Playlists: []*Playlist{&playlists[0], &playlists[1]}}}

	return Library{songs, playlists, folders}
}
