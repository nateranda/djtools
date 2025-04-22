# djtools
`djtools` is a library written in Go to manipulate and convert your DJ software libraries. It aims to be simple and fast, it depends on few third-party packages, and it's complete with a robust testing framework. Right now, the tool is only a Go library, but a CLI is planned for standalone use or to ship with other projects.

## Features

### Conversion
`djtools` converts your:
- Tracks
- Playlists and crates
- Cue points
- Hot cues
- Loops
- Beat grids

Smart playlist conversion and MP3 offset correction is planned.

### Platforms
Currently, `djtools` supports these platforms:
- Engine: import
- Rekordbox XML: export

`djtools` plans to support these platforms:
- Engine: export
- Rekordbox XML: import
- Rekordbox: import and export
- Algoriddim Djay: import and export
- Serato: import and export
- VirtualDJ: import and export
- Mixxx: import and export
- Spotify: playlist export
- Soundcloud: playlist export
- Beatport: playlist export

### Utilities
`djtools` plans to support these utilities:
- Purchase link finder: import a `djtools` library, a streaming service playlist, or a song name and get a list of links to buy the song(s) with their respective prices

## Usage
Below illustrates a basic usage of go
```go
package main

import (
	"log"

	"github.com/nateranda/djtools/engine"
	"github.com/nateranda/djtools/rbxml"
)

func main() {
  // initialize an empty options struct
  importOptions := engine.ImportOptions{}

  // import your Engine library
  library, err := engine.Import("import/path/", importOptions)
  if err != nil {
	  log.Panic(err)
  }

  // modify your library however you want,
  // like removing the first playlist
  library.Playlists = library.Playlists[1:]

  // export to a Rekordbox XML file
  err = rbxml.Export(&library, "export/path/library.xml")
  if err != nil {
	  log.Panic(err)
  }
}

```