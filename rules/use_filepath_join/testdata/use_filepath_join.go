package fixtures

import (
	"os"
	"path/filepath"
)

func badPathSeparatorConcat(dir, filename string) string {
	return dir + string(os.PathSeparator) + filename // MATCH /path concatenation using string(os.PathSeparator) can be replaced with filepath.Join/
}

func badPathSeparatorConcatMultiple(dir, sub, filename string) string {
	return dir + string(os.PathSeparator) + sub + string(os.PathSeparator) + filename // MATCH /path concatenation using string(os.PathSeparator) can be replaced with filepath.Join/
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

func goodSlashLiteralConcat(dir, filename string) string {
	return dir + "/" + filename // not flagged: literal slash is not string(os.PathSeparator)
}

func goodBackslashLiteralConcat(dir, filename string) string {
	return dir + "\\" + filename // not flagged: literal backslash is not string(os.PathSeparator)
}

func goodURLConcat(base, path string) string {
	return "https://example.com/" + path // not flagged: URL construction
}

func goodMQTTTopic(device string) string {
	return "devices/" + device + "/status" // not flagged: MQTT topic
}

func goodHTTPPath(version string) string {
	return "/api/" + version + "/users" // not flagged: HTTP path
}

func goodTrailingSlash(dir string) string {
	return dir + "/subdir" // not flagged: literal slash
}

var _ = os.PathSeparator // ensure os import is used
