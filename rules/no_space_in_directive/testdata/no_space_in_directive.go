// go:generate stringer -type=Pill
// MATCH:1 /compiler directive contains space: go:generate stringer -type=Pill/

package fixtures

// go:embed hello.txt
// MATCH:6 /compiler directive contains space: go:embed hello.txt/
var s string

// go:noinline
// MATCH:10 /compiler directive contains space: go:noinline/
func Add(a, b int) int {
	return a + b
}

// go:build linux
// MATCH:16 /compiler directive contains space: go:build linux/

// go:linkname localname importpath.name
// MATCH:19 /compiler directive contains space: go:linkname localname importpath.name/

// Valid directives below - no matches expected

//go:generate stringer -type=Pill

//go:embed hello.txt
var s2 string

//go:noinline
func Subtract(a, b int) int {
	return a - b
}

//go:build linux

//go:linkname localname2 importpath.name

// This is a normal comment mentioning go: but not a directive
// Just talking about go: things
