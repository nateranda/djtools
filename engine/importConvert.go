package engine

import (
	"encoding/binary"
	"math"
	"path/filepath"

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
		i += 8 // skip past unknown int32, not needed
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
		i += 8 // skip past unknown int32, not needed
		beatData.adjBeatgrid = append(beatData.adjBeatgrid, marker)
	}

	return beatData
}

func gridFromBeatData(sampleRate float64, enGrid []marker) []db.Marker {
	var grid []db.Marker
	for i := range len(enGrid) - 1 {
		var marker db.Marker
		marker.StartPosition = enGrid[i].offset / sampleRate
		lenMarker := enGrid[i+1].offset - enGrid[i].offset
		marker.Bpm = sampleRate * 60 * float64(enGrid[i].numBeats) / lenMarker
		marker.BeatNumber = int(enGrid[i].beatNumber) % 4
		grid = append(grid, marker)
	}

	return grid
}

func cuesFromBlob(sampleRate float64, blob []byte) cueData {
	var cueData cueData
	i := 8 //byte index, skipping number of cues (always 8)
	// skip unset cues
	for pos := range 8 {
		labelLength := int(blob[i])
		if labelLength == 0 { // label length 0 means no cue at this position
			i += 13
			continue
		}
		i++
		var cue db.HotCue
		cue.Position = pos + 1
		cue.Name = string(blob[i : i+labelLength])
		i += labelLength
		cue.Offset = math.Float64frombits(binary.BigEndian.Uint64(blob[i:i+8])) / sampleRate
		i += 8
		i++ // skip alpha channel (always 255)
		r := int(blob[i])
		i++
		g := int(blob[i])
		i++
		b := int(blob[i])
		i++
		color, err := db.RgbToHex(r, g, b)
		logError(err)
		cue.Color = color
		cueData.cues = append(cueData.cues, cue)
	}

	cueData.cueModified = math.Float64frombits(binary.BigEndian.Uint64(blob[i:i+8])) / sampleRate
	i += 9
	cueData.cueOriginal = math.Float64frombits(binary.BigEndian.Uint64(blob[i:i+8])) / sampleRate

	return cueData
}

func loopsFromBlob(sampleRate float64, blob []byte) []db.Loop {
	var loops []db.Loop
	i := 8 //byte index, skipping number of loops (always 8)
	for pos := range 8 {
		labelLength := int(blob[i])
		// skip unset loops
		if labelLength == 0 { // label length 0 means no loop at this position
			i += 23
			continue
		}
		i++
		var loop db.Loop
		loop.Position = pos + 1
		loop.Name = string(blob[i : i+labelLength])
		i += labelLength
		loop.Start = math.Float64frombits(binary.LittleEndian.Uint64(blob[i:i+8])) / sampleRate
		i += 8
		loop.End = math.Float64frombits(binary.LittleEndian.Uint64(blob[i:i+8])) / sampleRate
		i += 8
		i += 3 // skip alpha channel (always 255) and set bytes (not needed)
		r := int(blob[i])
		i++
		g := int(blob[i])
		i++
		b := int(blob[i])
		i++
		color, err := db.RgbToHex(r, g, b)
		logError(err)
		loop.Color = color
		loops = append(loops, loop)
	}

	return loops
}

func songHistoryFromHistoryList(historyList []historyListEntity) []songHistory {
	var songId int
	var lastPlayed int
	plays := 1

	var SongHistoryData []songHistory

	for i, HistoryListEntity := range historyList {
		if HistoryListEntity.id > songId && i != 0 {
			SongHistoryData = append(SongHistoryData, songHistory{songId, plays, lastPlayed})
			plays = 0
		}
		songId = HistoryListEntity.id
		lastPlayed = int(HistoryListEntity.startTime.Unix())
		plays += 1
	}
	SongHistoryData = append(SongHistoryData, songHistory{songId, plays, lastPlayed})

	return SongHistoryData
}

func fullPathFromRelativePath(basePath string, relativePath string) (string, error) {
	fullPath := filepath.Join(basePath, relativePath)
	absolutePath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}
	return absolutePath, nil
}

func importConvertSong(library *db.Library, songsNull []songNull, path string) {
	for _, song := range songsNull {
		path, err := fullPathFromRelativePath(path, song.path.String)
		logError(err)
		library.Songs = append(library.Songs, db.Song{
			SongID:       int(song.id.Int64),
			Title:        song.title.String,
			Artist:       song.artist.String,
			Composer:     song.composer.String,
			Album:        song.album.String,
			Genre:        song.genre.String,
			Filetype:     song.filetype.String,
			Size:         int(song.size.Int64),
			Length:       float32(song.length.Float64),
			Year:         int(song.year.Int64),
			Bpm:          float32(song.bpm.Float64),
			DateAdded:    int(song.dateAdded.Time.Unix()),
			DateModified: int(song.lastEditTime.Time.Unix()),
			Bitrate:      int(song.bitrate.Int64),
			Comment:      song.comment.String,
			Rating:       int(song.rating.Int64),
			Path:         path,
			Remixer:      song.remixer.String,
			Key:          song.key.String,
			Label:        song.label.String,
		})
	}
}

func importConvertPerformanceData(library *db.Library, perfData []performanceDataEntry) {
	for _, perfDataEntry := range perfData {
		song, err := db.GetSong(library.Songs, perfDataEntry.id)
		logError(err)

		beatDataBlob, err := qUncompress(perfDataEntry.beatDataBlob)
		logError(err)
		beatData := beatDataFromBlob(beatDataBlob)

		song.SampleRate = beatData.sampleRate

		song.Grid = gridFromBeatData(beatData.sampleRate, beatData.adjBeatgrid)

		quickCuesBlob, err := qUncompress(perfDataEntry.quickCuesBlob)
		logError(err)
		cueData := cuesFromBlob(beatData.sampleRate, quickCuesBlob)
		song.Cue = cueData.cueModified
		song.Cues = cueData.cues

		song.Loops = loopsFromBlob(beatData.sampleRate, perfDataEntry.loopsBlob)
	}
}

func importConvertHistory(library *db.Library, historyList []historyListEntity) {
	songHistoryData := songHistoryFromHistoryList(historyList)

	for _, songHistoryDataEntry := range songHistoryData {
		//fmt.Println(songHistoryDataEntry.id)
		song, err := db.GetSong(library.Songs, songHistoryDataEntry.id)
		// ignore any entries that don't match existing songs
		if err != nil {
			continue
		}
		song.PlayCount = songHistoryDataEntry.plays
		song.LastPlayed = songHistoryDataEntry.lastPlayed
	}
}

func importConvert(enLibrary library, path string) (db.Library, error) {
	var library db.Library
	importConvertSong(&library, enLibrary.songs, path)
	importConvertPerformanceData(&library, enLibrary.perfData)
	importConvertHistory(&library, enLibrary.historyList)
	return library, nil
}
