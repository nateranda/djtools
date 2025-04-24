package serato

type crate struct {
	filename string
	version  string
	paths    []string
}

type song struct {
	path  string
	geobs []geob
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
