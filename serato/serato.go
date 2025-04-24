package serato

func Import(path string) error {
	err := importExtract(path)
	if err != nil {
		return err
	}
	return nil
}
