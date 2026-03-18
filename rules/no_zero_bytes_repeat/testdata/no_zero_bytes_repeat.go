package fixtures

import "bytes"

func zeroBytesRepeatExamples() {
	// Invalid: count of 0 always returns empty slice
	padding := bytes.Repeat([]byte(" "), 0) // MATCH /bytes.Repeat called with a count of 0, always returns an empty slice/
	_ = padding

	// Valid: count is a variable
	n := 5
	padding2 := bytes.Repeat([]byte(" "), n)
	_ = padding2

	// Valid: count is a positive literal
	padding3 := bytes.Repeat([]byte("-"), 10)
	_ = padding3

	// Valid: empty slice created directly
	padding4 := []byte{}
	_ = padding4
}
