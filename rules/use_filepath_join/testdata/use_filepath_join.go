package fixtures

import "path/filepath"

func badSlashConcat(dir, filename string) string {
	return dir + "/" + filename // MATCH /path concatenation can be replaced with filepath.Join/
}

func badBackslashConcat(dir, filename string) string {
	return dir + "\\" + filename // MATCH /path concatenation can be replaced with filepath.Join/
}

func badTrailingSlash(dir string) string {
	return dir + "/subdir" // MATCH /path concatenation can be replaced with filepath.Join/
}

// Valid examples below: these should not trigger failures

func goodFilepathJoin(dir, filename string) string {
	return filepath.Join(dir, filename)
}

func goodStringConcat(a, b string) string {
	return a + b
}

func goodNonPathConcat(name string) string {
	return "hello " + name + "!"
}

func goodNumericAdd() int {
	return 1 + 2
}
