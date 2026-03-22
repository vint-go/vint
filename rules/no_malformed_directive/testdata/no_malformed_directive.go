package fixtures

// Invalid examples - should trigger failures

// go:noinline
// MATCH /malformed directive: space between // and go: in "// go:noinline"; use "//go:noinline"/
func spacedDirective() {}

//go:noinlin
// MATCH /malformed directive: possible misspelling of //go:noinline in "//go:noinlin"/
func misspelled() {}

//go:fakething
// MATCH /unrecognized directive "//go:fakething"/
func unrecognizedDirective() {}

//go:custom
// MATCH /unrecognized directive "//go:custom"/
func customDirective() {}

// Valid examples - should NOT trigger failures

//go:noinline
func validNoinline() {}

//go:generate stringer -type=MyType
func validGenerate() {}

//go:embed hello.txt
var embeddedFile string

// This is a normal comment, not a directive
func normalComment() {}
