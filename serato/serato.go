package serato

type crate struct {
	filename string
	version  string
	paths    []string
}

type song struct {
	path  string
	tags  tags
	geobs []geob
}

type tags struct {
	title       string
	album       string
	artist      string
	composer    string
	genre       string
	year        int
	trackNumber int
}

type geob struct {
	name  string
	value []byte
}

func Import(path string) error {
	err := importExtract(path)
	if err != nil {
		return err
	}
	return nil
}
