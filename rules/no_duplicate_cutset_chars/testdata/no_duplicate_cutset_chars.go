package fixtures

import (
	"fmt"
	"strings"
)

func badTrimLeftDuplicateChars() {
	// Likely intended TrimPrefix, not TrimLeft
	s := strings.TrimLeft("httpexample.com", "http") // MATCH /cutset "http" has duplicate characters, did you mean to use strings.TrimPrefix instead of strings.TrimLeft?/
	fmt.Println(s)
}

func badTrimRightDuplicateChars() {
	s := strings.TrimRight("example.comhttp", "http") // MATCH /cutset "http" has duplicate characters, did you mean to use strings.TrimSuffix instead of strings.TrimRight?/
	fmt.Println(s)
}

// Valid examples below: these should not trigger failures

func goodTrimPrefix() {
	// Use TrimPrefix to remove a prefix string
	s := strings.TrimPrefix("http://example.com", "http://")
	fmt.Println(s) // "example.com"
}

func goodTrimSuffix() {
	s := strings.TrimSuffix("example.com/path", "/path")
	fmt.Println(s)
}

func goodTrimLeftNoDuplicates() {
	// Cutset with all unique characters
	s := strings.TrimLeft("  hello", " ")
	fmt.Println(s)
}

func goodTrimRightNoDuplicates() {
	// Cutset with all unique characters
	s := strings.TrimRight("hello  ", " ")
	fmt.Println(s)
}

func goodTrimLeftUniqueChars() {
	s := strings.TrimLeft("abcdef", "abc")
	fmt.Println(s)
}
