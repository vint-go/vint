package fixtures

import "os"

// Invalid: deleting the entire temp directory with RemoveAll.
func badRemoveAll() {
	os.RemoveAll(os.TempDir()) // MATCH /deleting the system temp directory affects other processes; use a subdirectory instead/
}

// Invalid: deleting the temp directory with Remove.
func badRemove() {
	os.Remove(os.TempDir()) // MATCH /deleting the system temp directory affects other processes; use a subdirectory instead/
}

// Valid: removing a subdirectory created within the temp directory.
func goodRemoveSubdir() {
	dir, _ := os.MkdirTemp("", "myapp-")
	defer os.RemoveAll(dir)
}

// Valid: removing a specific path.
func goodRemoveSpecificPath() {
	os.RemoveAll("/tmp/myapp")
}

// Valid: removing a variable that is not os.TempDir().
func goodRemoveVariable() {
	path := "/some/path"
	os.RemoveAll(path)
}

// Valid: calling os.TempDir() alone without deletion.
func goodTempDirAlone() {
	_ = os.TempDir()
}
