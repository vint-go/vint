package fixtures

import (
	"bytes"
	"strings"
)

func redundantStringByteConversion() {
	var s string
	var b []byte

	// Invalid: bytes.Contains with []byte(s) conversion from string
	_ = bytes.Contains([]byte(s), []byte("pattern")) // MATCH /redundant []byte conversion in bytes.Contains call, use strings.Contains instead/

	// Invalid: bytes.HasPrefix with []byte conversion
	_ = bytes.HasPrefix([]byte(s), []byte("prefix")) // MATCH /redundant []byte conversion in bytes.HasPrefix call, use strings.HasPrefix instead/

	// Invalid: bytes.HasSuffix with []byte conversion
	_ = bytes.HasSuffix([]byte(s), []byte("suffix")) // MATCH /redundant []byte conversion in bytes.HasSuffix call, use strings.HasSuffix instead/

	// Invalid: bytes.Index with []byte conversion
	_ = bytes.Index([]byte(s), []byte("needle")) // MATCH /redundant []byte conversion in bytes.Index call, use strings.Index instead/

	// Invalid: bytes.Count with []byte conversion
	_ = bytes.Count([]byte(s), []byte("x")) // MATCH /redundant []byte conversion in bytes.Count call, use strings.Count instead/

	// Invalid: bytes.EqualFold with []byte conversion
	_ = bytes.EqualFold([]byte(s), []byte("test")) // MATCH /redundant []byte conversion in bytes.EqualFold call, use strings.EqualFold instead/

	// Invalid: strings.Contains with string(b) conversion from []byte
	_ = strings.Contains(string(b), "pattern") // MATCH /redundant string conversion in strings.Contains call, use bytes.Contains instead/

	// Invalid: strings.HasPrefix with string(b) conversion
	_ = strings.HasPrefix(string(b), "prefix") // MATCH /redundant string conversion in strings.HasPrefix call, use bytes.HasPrefix instead/

	// Invalid: strings.Index with string(b) conversion
	_ = strings.Index(string(b), "needle") // MATCH /redundant string conversion in strings.Index call, use bytes.Index instead/

	// Valid: bytes.Contains with actual []byte variable as first arg
	_ = bytes.Contains(b, []byte("pattern"))

	// Valid: strings.Contains with actual string variables
	_ = strings.Contains(s, "pattern")

	// Valid: strings.Contains with string literal
	_ = strings.Contains("hello world", "world")

	// Valid: bytes functions with no conversion on first arg
	_ = bytes.Index(b, []byte("needle"))

	// Valid: strings functions with no conversion
	_ = strings.Index(s, "needle")
}
