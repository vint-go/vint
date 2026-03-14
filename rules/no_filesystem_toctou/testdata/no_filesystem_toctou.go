package fixtures

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Invalid: os.ReadFile uses path from filepath.Walk callback
func processFilesWalk(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			data, err := os.ReadFile(path) // MATCH /potential TOCTOU race: os.ReadFile uses path from filepath.Walk callback/
			if err != nil {
				return err
			}
			_ = data
		}
		return nil
	})
}

// Invalid: os.Open uses path from filepath.WalkDir callback
func processFilesWalkDir(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			f, err := os.Open(path) // MATCH /potential TOCTOU race: os.Open uses path from filepath.WalkDir callback/
			if err != nil {
				return err
			}
			defer f.Close()
			_, _ = io.ReadAll(f)
		}
		return nil
	})
}

// Invalid: os.Stat uses path from filepath.Walk callback
func statInWalk(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		fi, err := os.Stat(path) // MATCH /potential TOCTOU race: os.Stat uses path from filepath.Walk callback/
		if err != nil {
			return err
		}
		_ = fi
		return nil
	})
}

// Invalid: os.Remove uses path from filepath.Walk callback
func removeInWalk(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err := os.Remove(path) // MATCH /potential TOCTOU race: os.Remove uses path from filepath.Walk callback/
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// Invalid: os.Chmod uses path from filepath.WalkDir callback
func chmodInWalkDir(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return os.Chmod(path, 0644) // MATCH /potential TOCTOU race: os.Chmod uses path from filepath.WalkDir callback/
	})
}

// Valid: no file operation using path in callback
func countFiles(root string) (int, error) {
	count := 0
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			count++
		}
		return nil
	})
	return count, err
}

// Valid: uses Fstat on file descriptor instead of path-based stat
func safeProcessFiles(root string) error {
	return filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			f, err := os.Open(path) // MATCH /potential TOCTOU race: os.Open uses path from filepath.WalkDir callback/
			if err != nil {
				return err
			}
			defer f.Close()

			fi, err := f.Stat()
			if err != nil {
				return err
			}
			if !fi.Mode().IsRegular() {
				return nil
			}
			data, err := io.ReadAll(f)
			if err != nil {
				return err
			}
			_ = data
		}
		return nil
	})
}

// Valid: file operation outside walk callback
func outsideWalk() {
	data, err := os.ReadFile("/etc/config")
	_ = data
	_ = err
}

// Valid: path variable is not the Walk callback's path parameter
func differentVariable(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		otherPath := "/some/hardcoded/path"
		data, err := os.ReadFile(otherPath)
		_ = data
		return err
	})
}
