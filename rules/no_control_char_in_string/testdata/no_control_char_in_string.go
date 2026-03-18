package fixtures

// Invalid: string with zero-width space (U+200B)
var s1 = "hello\u200bworld" // MATCH /string literal contains zero-width or control character U+200B/

// Invalid: string with zero-width joiner (U+200D)
var s2 = "test\u200dvalue" // MATCH /string literal contains zero-width or control character U+200D/

// Invalid: string with zero-width non-joiner (U+200C)
var s3 = "foo\u200cbar" // MATCH /string literal contains zero-width or control character U+200C/

// Invalid: string with soft hyphen (U+00AD)
var s4 = "soft\u00adhyphen" // MATCH /string literal contains zero-width or control character U+00AD/

// Invalid: string with left-to-right mark (U+200E)
var s5 = "ltr\u200emark" // MATCH /string literal contains zero-width or control character U+200E/

// Invalid: string with null byte (U+0000)
var s6 = "null\x00byte" // MATCH /string literal contains zero-width or control character U+0000/

// Valid: normal string
var v1 = "hello world"

// Valid: string with newline
var v2 = "hello\nworld"

// Valid: string with tab
var v3 = "hello\tworld"

// Valid: string with carriage return
var v4 = "hello\rworld"

// Valid: raw string
var v5 = `hello world`
