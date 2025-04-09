package db

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"errors"
	"io"
	"log"
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

type Cue struct {
	Name     string
	Offset   float64
	Position int
}

type Loop struct {
	Name     string
	Offset   float64
	Length   int
	Position int
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
	Grid         []Marker
	Cues         []Cue
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

func logError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// qUncompress uncompresses a uInt32-appended byte slice using zlib
// used for blobs compressed with the QT C++ library's qCompress function
func qUncompress(file []byte) ([]byte, error) {
	uncompressLength := binary.BigEndian.Uint32(file[:4])
	buffer := bytes.NewBuffer(file[4:])
	r, err := zlib.NewReader(buffer)
	logError(err)

	defer r.Close()

	var out bytes.Buffer
	_, err = io.Copy(&out, r)
	logError(err)

	fileDecomp := out.Bytes()

	// check if the file's uncompressed length matches the header
	if len(fileDecomp) != int(uncompressLength) {
		err := errors.New("db: uncompressed file length does not match length header")
		return []byte{}, err
	} else {
		return fileDecomp, nil
	}

}

// TBI
func Validate(library Library) (bool, string, error) {
	return true, "", nil
}

// TBI
func Import(program string, path string, options ImportOptions) (Library, error) {
	return Library{}, nil
}

// TBI
func Export(library Library, program string, path string, options ExportOptions) error {
	return nil
}
