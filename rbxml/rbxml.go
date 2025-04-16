package rbxml

import (
	"encoding/xml"
	"fmt"
	"os"

	"github.com/nateranda/djtools/db"
)

type product struct {
	Name    string `xml:"Name,attr"`
	Version string `xml:"Version,attr"`
	Company string `xml:"Company,attr"`
}

type track struct {
	TrackId      int     `xml:"TrackID,attr"`
	Name         string  `xml:"Name,attr"`
	Artist       string  `xml:"Artist,attr"`
	Composer     string  `xml:"Composer,attr"`
	Album        string  `xml:"Album,attr"`
	Grouping     string  `xml:"Grouping,attr"`
	Genre        string  `xml:"Genre,attr"`
	Kind         string  `xml:"Kind,attr"`
	Size         int64   `xml:"Size,attr"`
	TotalTime    float64 `xml:"TotalTime,attr"`
	DiscNumber   int32   `xml:"DiscNumber,attr"`
	TrackNumber  int32   `xml:"TrackNumber,attr"`
	Year         int32   `xml:"Year,attr"`
	AverageBpm   float64 `xml:"AverageBpm,attr"`
	DateModified string  `xml:"DateModidied,attr"` // yyyy-mm-dd
	DateAdded    string  `xml:"DateAdded,attr"`    // yyyy-mm-dd
	BitRate      int32   `xml:"BitRate,attr"`
	SampleRate   float64 `xml:"SampleRate,attr"`
	Comments     string  `xml:"Comments,attr"`
	PlayCount    int32   `xml:"PlayCount,attr"`
	LastPlayed   string  `xml:"LastPlayed,attr"` // yyyy-mm-dd
	Rating       int32   `xml:"Rating,attr"`
	Location     string  `xml:"Location,attr"` // URI formatted
	Remixer      string  `xml:"Remixer,attr"`
	Tonality     string  `xml:"Tonality,attr"`
	Label        string  `xml:"Label,attr"`
	Mix          string  `xml:"Mix,attr"`
	Colour       string  `xml:"Colour,attr"` // 0x-appended hex
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
	End      float64 `xml:"End,attr"`
	Num      int32   `xml:"Num,attr"` // hot cue: 0, 1, 2... memory cue: -1
}

type node struct {
	NodeType int32       `xml:"Type,attr"` // folder=0, playlist=1
	Name     string      `xml:"Name,attr"`
	Count    int32       `xml:"Count,attr"`   // number of sub-nodes
	Entries  int32       `xml:"Entries,attr"` // number of tracks in playlist
	KeyType  int32       `xml:"KeyType,attr"` // trackId=0, location=1, should always be 0
	Tracks   []nodeTrack `xml:"TRACK"`
	Nodes    []node      `xml:"NODE"`
}

type nodeTrack struct {
	Id int32 `xml:"Key,attr"`
}

type collection struct {
	Entries int32   `xml:"Entries,attr"` // number of tracks
	Tracks  []track `xml:"TRACK"`
}

type playlists struct {
	Nodes []node `xml:"NODE"`
}

type djPlaylists struct {
	XMLName    xml.Name   `xml:"DJPLAYLISTS"`
	Version    string     `xml:"Version,attr"` // should always be 1,0,0
	Product    product    `xml:"PRODUCT"`
	Collection collection `xml:"COLLECTION"`
	Playlists  playlists  `xml:"PLAYLISTS"`
}

func exportInsert(djPlaylists *djPlaylists, path string) error {
	xml, err := xml.MarshalIndent(djPlaylists, " ", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(path, xml, 0644)
	if err != nil {
		return fmt.Errorf("error exporting library: %v", err)
	}
	return nil
}

func Export(library *db.Library, path string) error {
	djPlaylists, err := exportConvert(library)
	if err != nil {
		return err
	}
	err = exportInsert(&djPlaylists, path)
	if err != nil {
		return err
	}
	return nil
}
