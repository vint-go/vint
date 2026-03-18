package fixtures

import "regexp"

// Invalid: package-level regexp.Compile with constant pattern
var re, _ = regexp.Compile(`\d+`) // MATCH /regexp.Compile can be replaced by regexp.MustCompile/

// Invalid: package-level regexp.CompilePOSIX with constant pattern
var rePosix, _ = regexp.CompilePOSIX(`[0-9]+`) // MATCH /regexp.CompilePOSIX can be replaced by regexp.MustCompilePOSIX/

// Valid: already using MustCompile
var reGood = regexp.MustCompile(`\d+`)

// Valid: already using MustCompilePOSIX
var reGoodPosix = regexp.MustCompilePOSIX(`[0-9]+`)

// Valid: pattern is not a string literal (variable)
var pattern = `\d+`
var reDynamic, _ = regexp.Compile(pattern)

func init() {
	// Invalid: regexp.Compile inside init function with constant pattern
	re, _ := regexp.Compile(`[a-z]+`) // MATCH /regexp.Compile can be replaced by regexp.MustCompile/
	_ = re
}

func normalFunc() {
	// Valid: regexp.Compile inside a regular function (not init, not package level)
	re, _ := regexp.Compile(`[a-z]+`)
	_ = re
}
