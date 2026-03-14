package fixtures

// Valid code: no bidirectional Unicode characters
func checkAccess(isAdmin bool) bool {
	if isAdmin {
		return true
	}
	return false
}

func greet(name string) string {
	return "Hello, " + name + "!"
}

// Line with RLO character: ‮ hidden
// MATCH /found bidirectional Unicode control character Right-to-Left Override (U+202E)/

// Line with LRO character: ‭ hidden
// MATCH /found bidirectional Unicode control character Left-to-Right Override (U+202D)/

// Line with RLE character: ‫ hidden
// MATCH /found bidirectional Unicode control character Right-to-Left Embedding (U+202B)/
