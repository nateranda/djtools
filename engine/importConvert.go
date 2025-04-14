package engine

import (
	"encoding/binary"
	"fmt"
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

func cuesFromBlob(sampleRate float64, blob []byte) (cueData, error) {
	var blobCueData cueData
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
		if err != nil {
			return cueData{}, fmt.Errorf("error extracting cues from cueData blob: %v", err)
		}
		cue.Color = color
		blobCueData.cues = append(blobCueData.cues, cue)
	}

	blobCueData.cueModified = math.Float64frombits(binary.BigEndian.Uint64(blob[i:i+8])) / sampleRate
	i += 9
	blobCueData.cueOriginal = math.Float64frombits(binary.BigEndian.Uint64(blob[i:i+8])) / sampleRate

	return blobCueData, nil
}

func loopsFromBlob(sampleRate float64, blob []byte) ([]db.Loop, error) {
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
		if err != nil {
			return nil, fmt.Errorf("error extracting loops from loops blob: %v", err)
		}
		loop.Color = color
		loops = append(loops, loop)
	}

	return loops, nil
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

func importConvertSong(library *db.Library, songsNull []songNull, path string) error {
	for _, song := range songsNull {
		songPath, err := fullPathFromRelativePath(path, song.path.String)
		if err != nil {
			return fmt.Errorf("error converting songs: %v", err)
		}
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
			Path:         songPath,
			Remixer:      song.remixer.String,
			Key:          song.key.String,
			Label:        song.label.String,
		})
	}
	return nil
}

func findFirstPlaylist(playlists []playlist) (int, error) {
	nextListIdSet := make(map[int]struct{})
	for _, playlist := range playlists {
		if playlist.nextListId != 0 {
			nextListIdSet[playlist.nextListId] = struct{}{}
		}
	}

	for _, playlist := range playlists {
		if _, exists := nextListIdSet[playlist.id]; !exists {
			return playlist.id, nil
		}
	}
	return 0, fmt.Errorf("NotFoundError: did not find the first playlist")
}

func sortPlaylists(playlists []playlist) ([]playlist, error) {
	i, err := findFirstPlaylist(playlists)
	if err != nil {
		return []playlist{}, fmt.Errorf("error sorting playlist: %v", err)
	}

	playlistMap := make(map[int]int)
	for i, playlist := range playlists {
		playlistMap[playlist.id] = i
	}
	var playlistsSorted []playlist
	for range playlists {
		playlistsSorted = append(playlistsSorted, playlists[playlistMap[i]])
		i = playlists[playlistMap[i]].nextListId
	}
	return playlistsSorted, nil
}

func populatePlaylists(library *db.Library, playlistEntityList []playlistEntity, playlists []playlist) []playlist {
	playlistMap := make(map[int]int)
	for i, playlist := range playlists {
		playlistMap[playlist.id] = i
	}

	songMap := make(map[int]int)
	for i, song := range library.Songs {
		songMap[song.SongID] = i
	}

	for _, playlistEntity := range playlistEntityList {
		trackId := playlistEntity.trackId
		listId := playlistEntity.listId
		playlists[playlistMap[listId]].songs = append(playlists[playlistMap[listId]].songs, &library.Songs[songMap[trackId]])
	}

	return playlists
}

func importConvertPerformanceData(library *db.Library, perfData []performanceDataEntry, importOptions ImportOptions) error {
	for _, perfDataEntry := range perfData {
		song, err := db.GetSong(library.Songs, perfDataEntry.id)
		if err != nil {
			return fmt.Errorf("error converting performance data: %v", err)
		}

		beatDataBlobComp := perfDataEntry.beatDataBlob
		if beatDataBlobComp == nil {
			fmt.Printf("Corrupt beatData blob for song id %d. Is the file corrupted? Skipping song...\n", perfDataEntry.id)
			continue
		}

		quickCuesBlobComp := perfDataEntry.quickCuesBlob
		if quickCuesBlobComp == nil {
			fmt.Printf("Corrupt quickCues blob for song id %d. Is the file corrupted? Skipping song...\n", perfDataEntry.id)
			continue
		}

		loopsBlob := perfDataEntry.loopsBlob
		if loopsBlob == nil {
			fmt.Printf("Corrupt loops blob for song id %d. Is the file corrupted? Skipping song...\n", perfDataEntry.id)
			continue
		}

		beatDataBlob, err := qUncompress(beatDataBlobComp)
		if err != nil {
			return err
		}
		beatData := beatDataFromBlob(beatDataBlob)

		song.SampleRate = beatData.sampleRate

		var beatgrid []marker
		if importOptions.ImportOriginalGrids {
			beatgrid = beatData.defaultBeatgrid
		} else {
			beatgrid = beatData.adjBeatgrid
		}

		song.Grid = gridFromBeatData(beatData.sampleRate, beatgrid)

		quickCuesBlob, err := qUncompress(quickCuesBlobComp)
		if err != nil {
			return fmt.Errorf("error converting performance data: %v", err)
		}
		cueData, err := cuesFromBlob(beatData.sampleRate, quickCuesBlob)
		if err != nil {
			return err
		}

		if importOptions.ImportOriginalCues {
			song.Cue = cueData.cueOriginal
		} else {
			song.Cue = cueData.cueModified
		}

		song.Cues = cueData.cues

		song.Loops, err = loopsFromBlob(beatData.sampleRate, loopsBlob)
		if err != nil {
			return err
		}
	}
	return nil
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

func importConvertPlaylist(library *db.Library, playlists []playlist, playlistEntityList []playlistEntity) error {
	playlists = populatePlaylists(library, playlistEntityList, playlists)

	parentPlaylistAddressMap := make(map[int]*db.Playlist)

	var parentPlaylists []playlist
	var playlistsNew []playlist

	// find parent playlists
	for _, playlist := range playlists {
		if playlist.parentListId == 0 {
			parentPlaylists = append(parentPlaylists, playlist)
		} else {
			playlistsNew = append(playlistsNew, playlist)
		}
	}
	// sort parent playlists
	parentPlaylists, err := sortPlaylists(parentPlaylists)
	if err != nil {
		return err
	}

	// move parent playlists to Library
	for i, playlist := range parentPlaylists {
		newPlaylist := db.Playlist{
			Name:       playlist.title,
			PlaylistID: playlist.id,
			Position:   i,
			Songs:      playlist.songs,
		}

		library.Playlists = append(library.Playlists, newPlaylist)

		// add new pointer(s) to map
		for i, playlist := range library.Playlists {
			parentPlaylistAddressMap[playlist.PlaylistID] = &library.Playlists[i]
		}
	}

	// remove parent playlists from playlists
	playlists = playlistsNew

	// iterate over child playlists until no child playlists left
	for len(playlists) > 0 {
		var parentPlaylistsNew []playlist
		var playlistsNew []playlist

		// generate map of parent playlists
		parentPlaylistIdMap := make(map[int]struct{})
		for _, playlist := range parentPlaylists {
			parentPlaylistIdMap[playlist.id] = struct{}{}
		}

		// find playlists whose parent was moved in the last round
		for _, playlist := range playlists {
			if _, exists := parentPlaylistIdMap[playlist.parentListId]; exists {
				parentPlaylistsNew = append(parentPlaylistsNew, playlist)
			} else {
				playlistsNew = append(playlistsNew, playlist)
			}
		}

		// sort new 'parent' playlists
		parentPlaylistsNew, err := sortPlaylists(parentPlaylistsNew)
		if err != nil {
			return err
		}

		// move new 'parent' playlists to library
		for i, playlist := range parentPlaylistsNew {
			newPlaylist := db.Playlist{
				Name:       playlist.title,
				PlaylistID: playlist.id,
				Position:   i,
				Songs:      playlist.songs,
			}
			parentPlaylist := parentPlaylistAddressMap[playlist.parentListId]
			parentPlaylist.SubPlaylists = append(parentPlaylist.SubPlaylists, newPlaylist)

			// add new pointer(s) to map
			for i, playlist := range parentPlaylist.SubPlaylists {
				parentPlaylistAddressMap[playlist.PlaylistID] = &parentPlaylist.SubPlaylists[i]
			}
		}

		// reset playlists and parentPlaylists with new slices
		playlists = playlistsNew
		parentPlaylists = parentPlaylistsNew
	}
	return nil
}

func importConvert(enLibrary library, path string, importOptions ImportOptions) (db.Library, error) {
	var library db.Library
	err := importConvertSong(&library, enLibrary.songs, path)
	if err != nil {
		return db.Library{}, err
	}
	err = importConvertPerformanceData(&library, enLibrary.perfData, importOptions)
	if err != nil {
		return db.Library{}, err
	}
	importConvertHistory(&library, enLibrary.historyList)
	err = importConvertPlaylist(&library, enLibrary.playlists, enLibrary.playlistEntityList)
	if err != nil {
		return db.Library{}, err
	}

	return library, nil
}
