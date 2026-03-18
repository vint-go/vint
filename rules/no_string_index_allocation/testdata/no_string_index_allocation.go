package fixtures

import (
	"bytes"
	"strings"
)

func stringIndexAllocation() {
	var b []byte

	// Invalid: strings.Index with byte-to-string conversion causes allocation
	_ = strings.Index(string(b), "pattern") // MATCH /strings.Index called with a byte-to-string conversion, use bytes.Index instead to avoid allocation/

	// Valid: using bytes.Index directly avoids the allocation
	_ = bytes.Index(b, []byte("pattern"))

	// Valid: strings.Index with an actual string variable
	s := "hello world"
	_ = strings.Index(s, "world")

	// Valid: strings.Index with a string literal
	_ = strings.Index("hello world", "world")
}
