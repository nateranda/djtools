package rbxml

import (
	"github.com/nateranda/djtools/db"
)

type product struct {
	name    string `xml:"Name,attr"`
	version string `xml:"Version,attr"`
	company string `xml:"Company,attr"`
}

type track struct {
	trackId      int     `xml:"TrackID,attr"`
	name         string  `xml:"Name,attr"`
	artist       string  `xml:"Artist,attr"`
	composer     string  `xml:"Composer,attr"`
	album        string  `xml:"Album,attr"`
	grouping     string  `xml:"Grouping,attr"`
	genre        string  `xml:"Genre,attr"`
	kind         string  `xml:"Kind,attr"`
	size         int64   `xml:"Size,attr"`
	totalTime    float64 `xml:"TotalTime,attr"`
	discNumber   int32   `xml:"DiscNumber,attr"`
	trackNumber  int32   `xml:"TrackNumber, attr"`
	year         int32   `xml:"Year,attr"`
	averageBpm   float64 `xml:"AverageBpm,attr"`
	dateModified string  `xml:"DateModidied,attr"` // yyyy-mm-dd
	dateAdded    string  `xml:"DateAdded,attr"`    // yyyy-mm-dd
	bitRate      int32   `xml:"BitRate,attr"`
	sampleRate   float64 `xml:"SampleRate,attr"`
	comments     string  `xml:"Comments,attr"`
	playCount    int32   `xml:"PlayCount,attr"`
	lastPlayed   string  `xml:"LastPlayed,attr"` // yyyy-mm-dd
	rating       int32   `xml:"Rating,attr"`
	location     string  `xml:"Location,attr"` // URI formatted
	remixer      string  `xml:"Remixer,attr"`
	tonality     string  `xml:"Tonality,attr"`
	label        string  `xml:"Label,attr"`
	mix          string  `xml:"Mix,attr"`
	colour       string  `xml:"Colour,attr"` // 0x-appended hex
}

type tempo struct {
	inizio  float64 `xml:"Inizio,attr"`
	bpm     float64 `xml:"Bpm,attr"`
	metro   string  `xml:"Metro,attr"` // 4/4, 3/4, etc.
	battito int32   `xml:"Battito,attr"`
}

type positionMark struct {
	name     string  `xml:"Name,attr"`
	markType int32   `xml:"Type,attr"` // cue=0, fade-in=1, fade-out=2, load=3, loop=4
	start    float64 `xml:"Start,attr"`
	end      float64 `xml:"End,attr"`
	num      int32   `xml:"Num,attr"` // hot cue: 0, 1, 2... memory cue: -1
}

type node struct {
	nodeType int32       `xml:"Type,attr"` // folder=0, playlist=1
	name     string      `xml:"Name,attr"`
	count    int32       `xml:"Count,attr"`   // number of sub-nodes
	entries  int32       `xml:"Entries,attr"` // number of tracks in playlist
	keyType  int32       `xml:"KeyType,attr"` // trackId=0, location=1, should always be 0
	tracks   []nodeTrack `xml:"TRACK"`
	nodes    []node      `xml:"NODE"`
}

type nodeTrack struct {
	id int32 `xml:"Key,attr"`
}

type collection struct {
	entries int32   `xml:"Entries,attr"` // number of tracks
	tracks  []track `xml:"TRACK"`
}

type playlists struct {
	nodes []node `xml:"NODE"`
}

type djPlaylists struct {
	version    string     `xml:"Version,attr"` // should always be 1,0,0
	product    product    `xml:"PRODUCT"`
	collection collection `xml:"COLLECTION"`
	playlists  playlists  `xml:"PLAYLISTS"`
}

func exportInsert(collection collection, path string) error {
	return nil
}

func Export(library *db.Library, path string) error {
	collection, err := exportConvert(library)
	if err != nil {
		return err
	}
	err = exportInsert(collection, path)
	if err != nil {
		return err
	}
	return nil
}
