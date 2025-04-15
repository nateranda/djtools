package db

// Marker is a marker in a beatgrid.
type Marker struct {
	StartPosition float64 // start position in seconds
	Bpm           float64 // beats per minute
	BeatNumber    int     // beat in measure to start on, 1-indexed, assumed 4/4
}

// HotCue is a saved hot cue for a song.
type HotCue struct {
	Name     string  // name of hot cue
	Offset   float64 // offset of hot cue in seconds
	Position int     // position of hot cue, 1-indexed
	Color    string  // color of hot cue button, hex code
}

// Loop is a saved loop for a song.
type Loop struct {
	Name     string  // name of loop
	Start    float64 // start time of loop, seconds
	End      float64 // end time of loop, seconds
	Position int     // position of loop, 1-indexed
	Color    string  // color of loop button, hex code
}

// Song is the metadata, analysis data, and saved cues/loops for a song.
type Song struct {
	SongID       int      // song id used by software
	Title        string   // title
	Artist       string   // artist
	Composer     string   // composer
	Album        string   // album song is from
	Grouping     string   // grouping
	Genre        string   // genre
	Filetype     string   // filetype, abbreviated
	Size         int      // file size in bytes
	Length       float32  // song length in seconds, approximate
	TrackNumber  int      // number in album
	Year         int      // release year
	Bpm          float32  // beats per minute
	DateModified int      // date last modified, unix
	DateAdded    int      // date added to library, unix
	Bitrate      int      // bitrate in kilobytes per second
	SampleRate   float64  // sample rate in hertz
	Comment      string   // comment
	PlayCount    int      // play count
	LastPlayed   int      // date last played, unix
	Rating       int      // rating in multiples of 20: 0*=0, 1*=20... 5*=100
	Path         string   // song absolute path
	Remixer      string   // remixer
	Key          string   // key in int representation of camelot, 0-indexed: 0=8B, 1=8A, 2=9B... 23=7A
	Label        string   // label
	Mix          string   // mix
	Color        string   // color in hex code
	Cue          float64  // cue locatin in seconds
	Grid         []Marker // slice of Marker structs in order
	Cues         []HotCue // slice of Cue structs out of order
	Loops        []Loop   // slice of Loop structs out of order
}

// Playlist is a playlist of songs which can contain other playlists.
// A folder is just a Playlist with no songs that contains other playlists.
type Playlist struct {
	PlaylistID   int        // playlist id used by software
	Name         string     // name of playlist
	Songs        []int      // slice of song ids in order
	SubPlaylists []Playlist // slice of child playlists in order, can be recursive
}

// Library is the entire library of a DJ software.
type Library struct {
	Songs     []Song     // slice of Song structs, can be ordered
	Playlists []Playlist // slice of Playlist structs in order
}
