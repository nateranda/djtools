package rbxml

import "github.com/nateranda/djtools/db"

func Export(library *db.Library, path string) error {
	err := exportConvert(library)
	if err != nil {
		return err
	}
	return nil
}
