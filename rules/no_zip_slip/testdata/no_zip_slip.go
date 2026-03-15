package fixtures

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Invalid: filepath.Join without strings.HasPrefix validation
func extractZipUnsafe(zipPath, dest string) error {
	r, _ := zip.OpenReader(zipPath)
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name) // MATCH /potential Zip Slip: archive entry path not validated after filepath.Join/
		outFile, _ := os.Create(path)
		rc, _ := f.Open()
		io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
	}
	return nil
}

// Valid: filepath.Join with strings.HasPrefix validation
func extractZipSafe(zipPath, dest string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()

	dest = filepath.Clean(dest) + string(os.PathSeparator)

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(path, dest) {
			return fmt.Errorf("illegal file path: %s", f.Name)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, 0750)
			continue
		}

		os.MkdirAll(filepath.Dir(path), 0750)
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
