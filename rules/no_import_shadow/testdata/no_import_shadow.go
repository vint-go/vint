package fixtures

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Invalid: shadows the filepath import
func invalidShadowFilepath() {
	filepath := "/some/path" // MATCH /assignment shadows import "filepath"/
	_ = filepath
}

// Invalid: shadows the fmt import
func invalidShadowFmt() {
	fmt := "hello" // MATCH /assignment shadows import "fmt"/
	_ = fmt
}

// Invalid: shadows the strings import
func invalidShadowStrings() {
	strings := []string{"a", "b"} // MATCH /assignment shadows import "strings"/
	_ = strings
}

// Valid: different variable name, does not shadow any import
func validNoShadow() {
	fp := "/some/path"
	fullPath := filepath.Join(fp, "file.txt")
	_ = fullPath
}

// Valid: uses the import normally
func validUseFmt() {
	msg := fmt.Sprintf("hello %s", "world")
	_ = msg
}

// Valid: different name
func validDifferentName() {
	str := strings.Join([]string{"a", "b"}, ",")
	_ = str
}
