package fixtures

import "regexp"

// Invalid: character class can be simplified to shorthand
var re1 = regexp.MustCompile(`[0-9]+`) // MATCH /regexp pattern can be simplified: `[0-9]+` -> `\d+`/

// Invalid: alternation can be factored with common prefix
var re2 = regexp.MustCompile(`http|https`) // MATCH /regexp pattern can be simplified: `http|https` -> `https?`/

// Invalid: negated character class can be simplified
var re3 = regexp.MustCompile(`[^0-9]`) // MATCH /regexp pattern can be simplified: `[^0-9]` -> `\D`/

// Invalid: POSIX class can be simplified
var re4 = regexp.MustCompile(`[[:space:]]`) // MATCH /regexp pattern can be simplified: `[[:space:]]` -> `\s`/

// Invalid: single char alternation can become char class
var re5 = regexp.MustCompile(`a|b|c`) // MATCH /regexp pattern can be simplified: `a|b|c` -> `[abc]`/

// Invalid: verbose quantifier {0,1} can be simplified to ?
var re6 = regexp.MustCompile(`x{0,1}`) // MATCH /regexp pattern can be simplified: `x{0,1}` -> `x?`/

// Invalid: verbose quantifier {1,} can be simplified to +
var re7 = regexp.MustCompile(`x{1,}`) // MATCH /regexp pattern can be simplified: `x{1,}` -> `x+`/

// Invalid: verbose quantifier {0,} can be simplified to *
var re8 = regexp.MustCompile(`x{0,}`) // MATCH /regexp pattern can be simplified: `x{0,}` -> `x*`/

// Invalid: common prefix factoring with multi-char suffix
var re9 = regexp.MustCompile(`foo|foobar`) // MATCH /regexp pattern can be simplified: `foo|foobar` -> `foo(?:bar)?`/

// Invalid: Compile also supported
var re10, _ = regexp.Compile(`[0-9]+`) // MATCH /regexp pattern can be simplified: `[0-9]+` -> `\d+`/

// Valid: already simplified
var reGood1 = regexp.MustCompile(`\d+`)

// Valid: already simplified
var reGood2 = regexp.MustCompile(`https?`)

// Valid: no simplification possible
var reGood3 = regexp.MustCompile(`[a-z]+`)

// Valid: pattern is not a string literal
var pattern = `[0-9]+`
var reDynamic = regexp.MustCompile(pattern)
