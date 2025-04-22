package engine

import (
	"encoding/binary"
	"fmt"
	"math"
	"path/filepath"

	"github.com/nateranda/djtools/lib"
)

func beatDataFromBlob(blob []byte) (beatData, error) {
	if len(blob) < 33 {
		return beatData{}, fmt.Errorf("InvalidBlobError: beatData blob should be at least 33 bytes")
	}

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
		i += 8 // skip unknown int32, not needed
		beatData.adjBeatgrid = append(beatData.adjBeatgrid, marker)
	}

	return beatData, nil
}

func gridFromBeatData(sampleRate float64, enGrid []marker) []lib.Marker {
	var grid []lib.Marker
	for i := range len(enGrid) - 1 {
		var marker lib.Marker
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
		var cue lib.HotCue
		cue.Position = pos + 1
		cue.Name = string(blob[i : i+labelLength])
		i += labelLength
		cue.Offset = math.Float64frombits(binary.BigEndian.Uint64(blob[i:i+8])) / sampleRate
		i += 9 // skip 1-byte alpha channel (always 255)
		r := int(blob[i])
		i++
		g := int(blob[i])
		i++
		b := int(blob[i])
		i++
		color, err := lib.RgbToHex(r, g, b)
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

func loopsFromBlob(sampleRate float64, blob []byte) ([]lib.Loop, error) {
	var loops []lib.Loop
	i := 8 // byte index, skipping number of loops (always 8)
	for pos := range 8 {
		labelLength := int(blob[i])
		// skip unset loops
		if labelLength == 0 { // label length 0 means no loop at this position
			i += 23
			continue
		}
		i++
		var loop lib.Loop
		loop.Position = pos + 1
		loop.Name = string(blob[i : i+labelLength])
		i += labelLength
		loop.Start = math.Float64frombits(binary.LittleEndian.Uint64(blob[i:i+8])) / sampleRate
		i += 8
		loop.End = math.Float64frombits(binary.LittleEndian.Uint64(blob[i:i+8])) / sampleRate
		i += 8
		i += 3 // skip 1-byte alpha channel (always 255) and set bytes (not needed)
		r := int(blob[i])
		i++
		g := int(blob[i])
		i++
		b := int(blob[i])
		i++
		color, err := lib.RgbToHex(r, g, b)
		if err != nil {
			return nil, fmt.Errorf("error extracting loops from loops blob: %v", err)
		}
		loop.Color = color
		loops = append(loops, loop)
	}

	return loops, nil
}

func fullPathFromRelativePath(basePath string, relativePath string) (string, error) {
	// hard to test - platform-specific
	fullPath := filepath.Join(basePath, relativePath)
	absolutePath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", err
	}
	return absolutePath, nil
}

func importConvertSong(library *lib.Library, songsNull []songNull, path string, importOptions ImportOptions) error {
	var err error
	for _, song := range songsNull {
		var songPath string
		if importOptions.PreserveOriginalPaths {
			songPath = song.path.String
		} else {
			songPath, err = fullPathFromRelativePath(path, song.path.String)
			if err != nil {
				return fmt.Errorf("error converting songs: %v", err)
			}
		}
		library.Songs = append(library.Songs, lib.Song{
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
			Key:          int(song.key.Int32),
			Label:        song.label.String,
		})
	}
	return nil
}

func findFirstPlaylist(playlists []playlist) (int, error) {
	nextListIdMap := make(map[int]struct{})
	for _, playlist := range playlists {
		if playlist.nextListId != 0 {
			nextListIdMap[playlist.nextListId] = struct{}{}
		}
	}

	// select the playlist that is not referred to by any other playlist's nextListId
	for _, playlist := range playlists {
		if _, exists := nextListIdMap[playlist.id]; !exists {
			return playlist.id, nil
		}
	}
	return 0, fmt.Errorf("NotFoundError: did not find the first playlist")
}

func findFirstSongs(playlistEntityList []playlistEntity) ([]playlistEntity, error) {
	var firstSongs []playlistEntity

	nextEntityIdMap := make(map[int]struct{})
	for _, playlistEntity := range playlistEntityList {
		if playlistEntity.nextEntityId != 0 {
			nextEntityIdMap[playlistEntity.nextEntityId] = struct{}{}
		}
	}

	// select only the songs that are not referred to by any other track's nextEntityId
	for _, playlistEntity := range playlistEntityList {
		if _, exists := nextEntityIdMap[playlistEntity.id]; !exists {
			firstSongs = append(firstSongs, playlistEntity)
		}
	}

	if firstSongs == nil {
		return nil, fmt.Errorf("NotFoundError: did not find any first songs")
	}

	return firstSongs, nil
}

func sortPlaylists(playlists []playlist) ([]playlist, error) {

	i, err := findFirstPlaylist(playlists)
	if err != nil {
		return nil, fmt.Errorf("error sorting playlist: %v", err)
	}

	playlistMap := make(map[int]int)
	for j, playlist := range playlists {
		playlistMap[playlist.id] = j
	}
	var playlistsSorted []playlist
	for range playlists {
		playlistsSorted = append(playlistsSorted, playlists[playlistMap[i]])
		i = playlists[playlistMap[i]].nextListId
	}
	return playlistsSorted, nil
}

func populatePlaylists(playlistEntityList []playlistEntity, playlists []playlist) ([]playlist, error) {
	playlistMap := make(map[int]int)
	for i, playlist := range playlists {
		playlistMap[playlist.id] = i
	}

	playlistEntityMap := make(map[int]*playlistEntity)
	for _, playlistEntity := range playlistEntityList {
		playlistEntityMap[playlistEntity.id] = &playlistEntity
	}

	firstSongs, err := findFirstSongs(playlistEntityList)
	if err != nil {
		return nil, fmt.Errorf("error populating playlists: %v", err)
	}

	for _, track := range firstSongs {
		// add first song
		trackId := track.trackId
		listId := track.listId
		nextEntityId := track.nextEntityId
		playlists[playlistMap[listId]].songs = append(playlists[playlistMap[listId]].songs, trackId)

		// iterate through playlist, adding songs in order until last song
		for range len(playlistEntityList) { // failsafe in case there is no last song
			if nextEntityId == 0 {
				break
			}
			trackId = playlistEntityMap[nextEntityId].trackId
			nextEntityId = playlistEntityMap[nextEntityId].nextEntityId
			playlists[playlistMap[listId]].songs = append(playlists[playlistMap[listId]].songs, trackId)
		}
	}

	return playlists, nil
}

func importConvertPerformanceData(library *lib.Library, perfData []performanceDataEntry, importOptions ImportOptions) error {
	songMap := make(map[int]*lib.Song)
	for i, song := range library.Songs {
		songMap[song.SongID] = &library.Songs[i]
	}

	for _, perfDataEntry := range perfData {
		song := songMap[perfDataEntry.id]

		beatDataBlobComp := perfDataEntry.beatDataBlob
		if beatDataBlobComp == nil {
			fmt.Printf("Corrupt beatData blob for song id %d. Marking song corrupt...\n", perfDataEntry.id)
			song.Corrupt = true
			continue
		}

		quickCuesBlobComp := perfDataEntry.quickCuesBlob
		loopsBlob := perfDataEntry.loopsBlob

		beatDataBlob, err := qUncompress(beatDataBlobComp)
		if err != nil {
			return err
		}
		beatData, err := beatDataFromBlob(beatDataBlob)
		if err != nil {
			return err
		}

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

func importConvertHistory(library *lib.Library, songHistoryList []songHistory) {
	songMap := make(map[int]*lib.Song)
	for i, song := range library.Songs {
		songMap[song.SongID] = &library.Songs[i]
	}

	for _, songHistory := range songHistoryList {
		song := songMap[songHistory.id]
		// ignore any entries for removed songs
		if song == nil {
			continue
		}
		song.PlayCount = songHistory.plays
		song.LastPlayed = songHistory.lastPlayed
	}
}

func importConvertPlaylist(library *lib.Library, playlists []playlist, playlistEntityList []playlistEntity) error {
	if playlists == nil {
		return nil
	}

	playlists, err := populatePlaylists(playlistEntityList, playlists)
	if err != nil {
		return err
	}

	parentPlaylistAddressMap := make(map[int]*lib.Playlist)

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
	parentPlaylists, err = sortPlaylists(parentPlaylists)
	if err != nil {
		return err
	}

	// move parent playlists to Library
	for _, playlist := range parentPlaylists {
		newPlaylist := lib.Playlist{
			Name:       playlist.title,
			PlaylistID: playlist.id,
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
		for _, playlist := range parentPlaylistsNew {
			newPlaylist := lib.Playlist{
				Name:       playlist.title,
				PlaylistID: playlist.id,
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

func removeSongFromPlaylists(playlists []lib.Playlist, songID int) []lib.Playlist {
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

func importCheckCorruptedSongs(library *lib.Library) {
	for i, song := range library.Songs {
		// this is expensive, but it should happen rarely so it's ok
		if song.Corrupt {
			// remove song from library.Songs (doesn't preserve order)
			library.Songs[i] = library.Songs[len(library.Songs)-1]
			library.Songs = library.Songs[:len(library.Songs)-1]

			// remove song from playlists
			library.Playlists = removeSongFromPlaylists(library.Playlists, song.SongID)
		}
	}
}

func importConvert(enLibrary library, path string, importOptions ImportOptions) (lib.Library, error) {
	var library lib.Library
	err := importConvertSong(&library, enLibrary.songs, path, importOptions)
	if err != nil {
		return lib.Library{}, err
	}
	err = importConvertPerformanceData(&library, enLibrary.perfData, importOptions)
	if err != nil {
		return lib.Library{}, err
	}
	importConvertHistory(&library, enLibrary.songHistoryList)
	err = importConvertPlaylist(&library, enLibrary.playlists, enLibrary.playlistEntityList)
	if err != nil {
		return lib.Library{}, err
	}

	importCheckCorruptedSongs(&library)

	return library, nil
}
