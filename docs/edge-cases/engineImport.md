# Edge Cases: Engine Import

## `engine.go`

### `qUncompress`
* *byte slice is at least 5 bytes long -> verify byte slice length*

### `initDB`
no assumptions

## `importExtract.go`

### `importExtractTrack`
* database pointer points to `m.db` -> instead, pass in databases struct and choose m.db
* table Track is in `m.db` -> OK, throws error if not
* table entries are in the right datatypes -> OK, throws error if not
* *every entry has a trackId -> filter for entries with trackId in query*

### `importExtractHistory`
* database pointer points to `hm.db` -> OK, tables are specific to `hm.db`
* tables Track and HistoryListEntity are in `hm.db` -> OK, throws error if not
* table entries are in the right datatypes -> OK, throws error if not
* *every entry has an originTrackId and a startTime -> filter for entries with both an originTrackId and a startTime in query*

### `importExtractPerformanceData`
* *database pointer points to `m.db` -> instead, pass in databases struct and choose m.db*
* table PerformanceData is in `m.db` -> OK, throws error if not
* table entries are in the right datatypes -> OK, throws error if not
* *every entry has a trackId -> filter for entries with trackId in query*

### `importExtractPlaylist`
* *database pointer points to `m.db` -> instead, pass in databases struct and choose m.db*
* table Playlist is in `m.db` -> OK, throws error if not
* table entries are in the right datatypes -> OK, throws error if not
* *every entry has an id -> filter for entries with id in query*

### `importExtractPlaylistEntity`
* *database pointer points to `m.db` -> instead, pass in databases struct and choose m.db*
* table PlaylistEntity is in `m.db` -> OK, throws error if not
* table entries are in the right datatypes -> OK, throws error if not
* *every entry has a listId -> filter for entries with listId in query*

### `importExtractSmartlist`
* *database pointer points to `m.db` -> instead, pass in databases struct and choose m.db*
* table Smartlist is in `m.db` -> OK, throws error if not
* table entries are in the right datatypes -> OK, throws error if not
* *every entry has a listUuid -> filter for entries with listUuid in query*

### `importExtract`
* *initDB returns `m.db` first and `hm.db` second -> instead, return struct with values `m` and `hm` pointing to databases*

## `importConvert.go`

### `beatDataFromBlob`
* *blob is at least 33 bytes long -> verify byte slice length*
* *numMarkers is a normal number -> verify numMarkers isn't more than 16 or so*

### `gridFromBeatData`
* *sampleRate is positive and not 0 -> verify sampleRate*
* *lenMarker is positive and not 0 -> verify lenMarker*

### `cuesFromBlob`
* number of cues is always 8 -> OK, to spec
* each unset cue is 13 bytes long -> OK, to spec
* *blob is at least 129 bytes long -> verify byte slice length*
* *sampleRate is positive and not 0 -> verify sampleRate*

### `loopsFromBlob`
* number of loops is always 8 -> OK, to spec
* each unset loop is 23 bytes long -> OK, to spec
* *blob is at least 192 bytes long -> verify byte slice length*

### `songHistoryFromHistoryList`
* *assumes historyList is at least 1 entry long -> return nil if empty*

### `fullPathFromRelativePath`
* *basePath and relativePath are valid paths -> validate paths*

### `importConvertSong`
* songPath is valid -> OK, validated in fullPathFromRelativePath

### `findFirstPlaylist`
* *playlists is not empty -> verify and return 0 if empty*
* *playlists have unique ids -> verify and throw error if duplicate*

### `findFirstSongs`
* *playlistEntityList is not empty -> verify and return nil if empty*

### `sortPlaylists`
* *playlists is not empty -> verify and return nil if empty*
* *playlists have nextListIds -> verify and throw error if empty*

### `populatePlaylists`
* *playlists have unique ids -> verify and throw error if duplicate*
* *playlistEntities have unique ids -> verify and throw error if duplicate*
* *playlistEntityList is not empty -> verify and return nil if empty*
* *nextEntityIds point to actual ids -> verify and throw error if false*

### `importConvertPerformanceData`
* 