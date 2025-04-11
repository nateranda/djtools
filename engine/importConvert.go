package engine

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"

	"github.com/nateranda/djtools/db"
)

func beatDataFromBlob(blob []byte) beatData {
	var beatData beatData
	i := 0 // byte index

	// get sample rate
	beatData.sampleRate = math.Float64frombits(binary.BigEndian.Uint64(blob[i : i+8]))
	i += 17 // skip past track length and beatgrid set byte, not needed

	// save normal beatgrid
	numMarkers := binary.BigEndian.Uint64(blob[i : i+8])
	i += 8

	for range numMarkers {
		var marker marker
		marker.offset = math.Float64frombits(binary.LittleEndian.Uint64(blob[i : i+8]))
		i += 8
		marker.beatNumber = int64(binary.LittleEndian.Uint64(blob[i : i+8]))
		i += 8
		marker.numBeats = binary.LittleEndian.Uint32(blob[i : i+4])
		i += 8
		beatData.defaultBeatgrid = append(beatData.defaultBeatgrid, marker)
	}

	// save adjusted beatgrid
	numMarkers = binary.BigEndian.Uint64(blob[i : i+8])
	i += 8
	for range numMarkers {
		var marker marker
		marker.offset = math.Float64frombits(binary.LittleEndian.Uint64(blob[i : i+8]))
		i += 8
		marker.beatNumber = int64(binary.LittleEndian.Uint64(blob[i : i+8]))
		i += 8
		marker.numBeats = binary.LittleEndian.Uint32(blob[i : i+4])
		i += 8
		beatData.adjBeatgrid = append(beatData.adjBeatgrid, marker)
	}

	return beatData
}

// unused
func importConvertSong(song SongNull) db.Song {
	return db.Song{
		SongID:    int(song.SongID.Int64),
		Title:     song.Title.String,
		Artist:    song.Artist.String,
		Composer:  song.Composer.String,
		Album:     song.Album.String,
		Genre:     song.Genre.String,
		Filetype:  song.Filetype.String,
		Size:      int(song.Size.Int64),
		Length:    float32(song.Length.Float64),
		Year:      int(song.Year.Int64),
		Bpm:       float32(song.Bpm.Float64),
		DateAdded: int(song.DateAdded.Time.Unix()),
		Bitrate:   int(song.Bitrate.Int64),
		Comment:   song.Comment.String,
		Rating:    int(song.Rating.Int64),
		Path:      song.Path.String,
		Remixer:   song.Remixer.String,
		Key:       song.Key.String,
		Label:     song.Label.String,
	}
}

// unused
func importConvertSongHistory(historyList []HistoryListEntity) []SongHistory {
	var songId int
	var lastPlayed int
	plays := 1

	var SongHistoryData []SongHistory

	for i, HistoryListEntity := range historyList {
		if HistoryListEntity.trackId > songId && i != 0 {
			SongHistoryData = append(SongHistoryData, SongHistory{songId, plays, lastPlayed})
			plays = 0
		}
		songId = HistoryListEntity.trackId
		lastPlayed = int(HistoryListEntity.startTime.Unix())
		plays += 1
	}
	SongHistoryData = append(SongHistoryData, SongHistory{songId, plays, lastPlayed})

	return SongHistoryData
}

func ImportConvertGrid() {
	beatDataComp, err := os.ReadFile("tmp/beatData")
	logError(err)

	beatDataBlob, err := qUncompress(beatDataComp)
	logError(err)

	fmt.Println(beatDataBlob)

	beatData := beatDataFromBlob(beatDataBlob)
	fmt.Println(beatData)
}
