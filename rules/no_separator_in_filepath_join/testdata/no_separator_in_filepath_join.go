package fixtures

import "path/filepath"

func badTrailingSlash(filename string) string {
	return filepath.Join("dir/", filename) // MATCH /filepath.Join argument "dir/" contains a path separator/
}

func badBackslash(filename string) string {
	return filepath.Join("base\\sub", filename) // MATCH /filepath.Join argument "base\\sub" contains a path separator/
}

func badEmbeddedSlash(filename string) string {
	return filepath.Join("a/b", "c", filename) // MATCH /filepath.Join argument "a/b" contains a path separator/
}

// Valid examples below: these should not trigger failures

func goodNoSeparator(filename string) string {
	return filepath.Join("dir", filename)
}

func goodMultipleArgs(filename string) string {
	return filepath.Join("base", "sub", filename)
}

func goodSingleArg() string {
	return filepath.Join("file.txt")
}

func goodVariable(dir, filename string) string {
	return filepath.Join(dir, filename)
}
