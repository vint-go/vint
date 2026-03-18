package fixtures

import "strings"

// Invalid: passing invalid UTF-8 string to strings.ToUpper via variable.
func badToUpperViaVar() {
	s := string([]byte{0xff, 0xfe})
	result := strings.ToUpper(s) // MATCH /argument to strings.ToUpper is not a valid UTF-8 encoded string/
	_ = result
}

// Invalid: passing invalid UTF-8 string directly to strings.ToLower.
func badToLowerDirect() {
	result := strings.ToLower(string([]byte{0xff, 0xfe})) // MATCH /argument to strings.ToLower is not a valid UTF-8 encoded string/
	_ = result
}

// Invalid: passing invalid UTF-8 string to strings.Contains.
func badContains() {
	s := string([]byte{0x80, 0x81})
	_ = strings.Contains("hello", s) // MATCH /argument to strings.Contains is not a valid UTF-8 encoded string/
}

// Invalid: passing invalid UTF-8 string to strings.TrimLeft.
func badTrimLeft() {
	s := string([]byte{0xc0, 0xaf})
	_ = strings.TrimLeft("hello", s) // MATCH /argument to strings.TrimLeft is not a valid UTF-8 encoded string/
}

// Valid: passing a valid UTF-8 string to strings.ToUpper.
func goodToUpper() {
	s := "hello world"
	result := strings.ToUpper(s)
	_ = result
}

// Valid: passing a valid UTF-8 string via byte conversion.
func goodValidBytes() {
	s := string([]byte{0x48, 0x65, 0x6c, 0x6c, 0x6f}) // "Hello"
	result := strings.ToUpper(s)
	_ = result
}

// Valid: no strings function call.
func goodNoStringsCall() {
	s := string([]byte{0xff, 0xfe})
	_ = s
}

// Valid: strings function with string literal.
func goodStringLiteral() {
	result := strings.ToUpper("hello")
	_ = result
}
