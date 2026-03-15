package fixtures

import (
	"bytes"
	"strings"
)

func noopStringsReplace(s string) {
	strings.Replace(s, "old", "new", 0) // MATCH /strings.Replace called with n == 0: replaces 0 occurrences (no-op)/
}

func noopStringsSplitN(s string) {
	strings.SplitN(s, ",", 0) // MATCH /strings.SplitN called with n == 0: returns nil/
}

func noopStringsSplitAfterN(s string) {
	strings.SplitAfterN(s, ",", 0) // MATCH /strings.SplitAfterN called with n == 0: returns nil/
}

func noopBytesReplace(b []byte) {
	bytes.Replace(b, []byte("old"), []byte("new"), 0) // MATCH /bytes.Replace called with n == 0: replaces 0 occurrences (no-op)/
}

func noopBytesSplitN(b []byte) {
	bytes.SplitN(b, []byte(","), 0) // MATCH /bytes.SplitN called with n == 0: returns nil/
}

func noopBytesSplitAfterN(b []byte) {
	bytes.SplitAfterN(b, []byte(","), 0) // MATCH /bytes.SplitAfterN called with n == 0: returns nil/
}

// Valid examples below: these should not trigger failures

func validStringsReplace(s string) string {
	return strings.Replace(s, "old", "new", -1)
}

func validStringsReplaceAll(s string) string {
	return strings.ReplaceAll(s, "old", "new")
}

func validStringsSplitN(s string) []string {
	return strings.SplitN(s, ",", -1)
}

func validStringsSplitNPositive(s string) []string {
	return strings.SplitN(s, ",", 3)
}

func validStringsReplacePositive(s string) string {
	return strings.Replace(s, "old", "new", 1)
}

func validBytesReplace(b []byte) []byte {
	return bytes.Replace(b, []byte("old"), []byte("new"), -1)
}

func validBytesSplitN(b []byte) [][]byte {
	return bytes.SplitN(b, []byte(","), 2)
}
