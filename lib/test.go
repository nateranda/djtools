package lib

import (
	"fmt"
	"io"
	"os"
	"testing"
)

// CopyFile copies a file from a source to a destination
func CopyFile(t *testing.T, srcPath string, destPath string) {
	t.Helper()
	srcFile, err := os.Open(srcPath)
	if err != nil {
		t.Errorf("unexpected error opening source file: %v", err)
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		t.Errorf("unexpected error creating destination file: %v", err)
	}

	i, err := io.Copy(destFile, srcFile)
	if err != nil {
		t.Errorf("unexpected error copying source file: %v", err)
	}
	fmt.Println(i)

	err = destFile.Sync()
	if err != nil {
		t.Errorf("unexpected error syncing destination file: %v", err)
	}
}
