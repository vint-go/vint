package fixtures

import "fmt"

func usePercentQInvalid() {
	s := "hello"
	_ = fmt.Sprintf("value is \"%s\"", s) // MATCH /use %q instead of \"%s\" for quoted string formatting/

	_ = fmt.Sprintf("key \"%s\" not found", s) // MATCH /use %q instead of \"%s\" for quoted string formatting/

	_ = fmt.Printf("got \"%s\" from input", s) // MATCH /use %q instead of \"%s\" for quoted string formatting/

	_ = fmt.Errorf("unknown field \"%s\"", s) // MATCH /use %q instead of \"%s\" for quoted string formatting/
}

func usePercentQValid() {
	s := "hello"
	// Using %q is correct
	_ = fmt.Sprintf("value is %q", s)

	// No surrounding quotes around %s
	_ = fmt.Sprintf("value is %s", s)

	// Using %s without escaped quotes
	_ = fmt.Sprintf("value: %s", s)

	// Different verb, not %s
	_ = fmt.Sprintf("value is \"%d\"", 42)

	// Multiple args without the pattern
	_ = fmt.Sprintf("%s = %s", "a", "b")
}
