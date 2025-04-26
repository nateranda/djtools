// This package contains import and export functions for Rekordbox's XML format.
package rbxml

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"

	"github.com/nateranda/djtools/lib"
)

const version string = "0.1"

type ExportOptions struct {
	UseUTC bool
}

type product struct {
	Name    string `xml:"Name,attr"`
	Version string `xml:"Version,attr"`
	Company string `xml:"Company,attr"`
}

type track struct {
	TrackId      int             `xml:"TrackID,attr"`
	Name         string          `xml:"Name,attr,omitempty"`
	Artist       string          `xml:"Artist,attr,omitempty"`
	Composer     string          `xml:"Composer,attr,omitempty"`
	Album        string          `xml:"Album,attr,omitempty"`
	Grouping     string          `xml:"Grouping,attr,omitempty"`
	Genre        string          `xml:"Genre,attr,omitempty"`
	Kind         string          `xml:"Kind,attr,omitempty"`
	Size         int64           `xml:"Size,attr,omitempty"`
	TotalTime    float64         `xml:"TotalTime,attr,omitempty"`
	DiscNumber   int32           `xml:"DiscNumber,attr,omitempty"`
	TrackNumber  int32           `xml:"TrackNumber,attr,omitempty"`
	Year         int32           `xml:"Year,attr,omitempty"`
	AverageBpm   float64         `xml:"AverageBpm,attr,omitempty"`
	DateModified string          `xml:"DateModidied,attr,omitempty"` // yyyy-mm-dd
	DateAdded    string          `xml:"DateAdded,attr,omitempty"`    // yyyy-mm-dd
	BitRate      int32           `xml:"BitRate,attr,omitempty"`
	SampleRate   float64         `xml:"SampleRate,attr,omitempty"`
	Comments     string          `xml:"Comments,attr,omitempty"`
	PlayCount    int32           `xml:"PlayCount,attr,omitempty"`
	LastPlayed   string          `xml:"LastPlayed,attr,omitempty"` // yyyy-mm-dd
	Rating       int32           `xml:"Rating,attr,omitempty"`
	Location     string          `xml:"Location,attr,omitempty"` // URI formatted
	Remixer      string          `xml:"Remixer,attr,omitempty"`
	Tonality     string          `xml:"Tonality,attr,omitempty"`
	Label        string          `xml:"Label,attr,omitempty"`
	Mix          string          `xml:"Mix,attr,omitempty"`
	Colour       string          `xml:"Colour,attr,omitempty"` // 0x-appended hex
	Tempo        *[]tempo        `xml:"TEMPO,omitempty"`
	PositionMark *[]positionMark `xml:"POSITION_MARK,omitempty"`
}

type tempo struct {
	Inizio  float64 `xml:"Inizio,attr"`
	Bpm     float64 `xml:"Bpm,attr"`
	Metro   string  `xml:"Metro,attr"` // 4/4, 3/4, etc.
	Battito int32   `xml:"Battito,attr"`
}

type positionMark struct {
	Name     string  `xml:"Name,attr"`
	MarkType int32   `xml:"Type,attr"` // cue=0, fade-in=1, fade-out=2, load=3, loop=4
	Start    float64 `xml:"Start,attr"`
	End      float64 `xml:"End,attr,omitempty"`
	Num      int32   `xml:"Num,attr"` // hot cue: 0, 1, 2... memory cue: -1
}

type node struct {
	NodeType int32        `xml:"Type,attr"` // folder=0, playlist=1
	Name     string       `xml:"Name,attr"`
	Count    int32        `xml:"Count,attr,omitempty"`   // number of sub-nodes
	Entries  int32        `xml:"Entries,attr,omitempty"` // number of tracks in playlist
	KeyType  int32        `xml:"KeyType,attr"`           // trackId=0, location=1, should always be 0
	Tracks   *[]nodeTrack `xml:"TRACK,omitempty"`
	Nodes    *[]node      `xml:"NODE,omitempty"`
}

type nodeTrack struct {
	Id int32 `xml:"Key,attr"`
}

type collection struct {
	Entries int32   `xml:"Entries,attr"` // number of tracks
	Tracks  []track `xml:"TRACK"`
}

type playlists struct {
	Node node `xml:"NODE"`
}

type djPlaylists struct {
	XMLName    xml.Name   `xml:"DJ_PLAYLISTS"`
	Version    string     `xml:"Version,attr"` // should always be 1,0,0
	Product    product    `xml:"PRODUCT"`
	Collection collection `xml:"COLLECTION"`
	Playlists  playlists  `xml:"PLAYLISTS"`
}

// sort sorts a djPlaylists struct's songs by song id, then each
// song's PositionMarks by cue, then cue points by id, then loops by id.
// This is to standardize XML output for testing
func (d *djPlaylists) sort() {
	sort.Slice(d.Collection.Tracks, func(i, j int) bool {
		return d.Collection.Tracks[i].TrackId < d.Collection.Tracks[j].TrackId
	})

	for i, track := range d.Collection.Tracks {
		positionMarks := *track.PositionMark
		sort.Slice(positionMarks, func(i, j int) bool {
			// Sort by MarkType first
			if positionMarks[i].MarkType != positionMarks[j].MarkType {
				return positionMarks[i].MarkType < positionMarks[j].MarkType
			}
			// If MarkType is the same, sort by Num
			return positionMarks[i].Num < positionMarks[j].Num
		})
		d.Collection.Tracks[i].PositionMark = &positionMarks
	}
}

// write writes a djPlaylists struct to a XML file at the given path
func (d *djPlaylists) write(path string) error {
	xml, err := xml.MarshalIndent(d, " ", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, xml, 0644)
	if err != nil {
		return fmt.Errorf("error exporting library: %v", err)
	}
	return nil
}

// read reads a djPlaylists struct from a XML file at the given path
func (d *djPlaylists) read(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	err = xml.Unmarshal(data, d)
	if err != nil {
		return fmt.Errorf("error unmarshaling XML file: %v", err)
	}

	return nil
}

func Import(path string) (lib.Library, error) {
	var djPlaylists djPlaylists
	err := djPlaylists.read(path)
	if err != nil {
		return lib.Library{}, err
	}

	library, err := importConvert(&djPlaylists)
	if err != nil {
		return lib.Library{}, err
	}

	library.CheckCorruptedSongs()

	return library, nil
}

func Export(library *lib.Library, path string, options ExportOptions) error {
	djPlaylists, err := exportConvert(library, options)
	if err != nil {
		return err
	}
	djPlaylists.sort() // for test parity
	err = djPlaylists.write(path)
	if err != nil {
		return err
	}
	return nil
}
