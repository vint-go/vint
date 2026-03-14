package fixtures

// Invalid directives - should trigger failures

//go:embod hello.txt
// MATCH /compiler directive unrecognized: embod/
var s string

//go:genrate stringer -type=Pill
// MATCH /compiler directive unrecognized: genrate/

//go:linknme localname importpath.name
// MATCH /compiler directive unrecognized: linknme/
func localname()

//go:noInline
// MATCH /compiler directive unrecognized: noInline/
func Add(a, b int) int {
	return a + b
}

//go:inline
// MATCH /compiler directive unrecognized: inline/
func Multiply(a, b int) int {
	return a * b
}

// Valid directives - should NOT trigger failures

//go:embed hello.txt
var s2 string

//go:generate stringer -type=Pill

//go:linkname localname2 importpath.name

//go:noinline
func Subtract(a, b int) int {
	return a - b
}

//go:build linux

//go:nosplit
func criticalFunc() {
}
