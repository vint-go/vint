package fixtures

import "strings"

func invalid(s string) {
	if strings.Index(s, "hello") != -1 { // MATCH /replace call to strings.Index with strings.Contains/
		// found
	}

	if strings.Index(s, "hello") >= 0 { // MATCH /replace call to strings.Index with strings.Contains/
		// found
	}

	if strings.Index(s, "hello") > -1 { // MATCH /replace call to strings.Index with strings.Contains/
		// found
	}

	if strings.Index(s, "hello") == -1 { // MATCH /replace call to strings.Index with strings.Contains/
		// not found
	}

	if strings.Index(s, "hello") < 0 { // MATCH /replace call to strings.Index with strings.Contains/
		// not found
	}

	// Reversed operand order
	if -1 != strings.Index(s, "world") { // MATCH /replace call to strings.Index with strings.Contains/
		// found
	}

	if 0 <= strings.Index(s, "world") { // MATCH /replace call to strings.Index with strings.Contains/
		// found
	}

	if -1 < strings.Index(s, "world") { // MATCH /replace call to strings.Index with strings.Contains/
		// found
	}

	if -1 == strings.Index(s, "world") { // MATCH /replace call to strings.Index with strings.Contains/
		// not found
	}

	if 0 > strings.Index(s, "world") { // MATCH /replace call to strings.Index with strings.Contains/
		// not found
	}
}

func valid(s string) {
	if strings.Contains(s, "hello") {
		// correct usage
	}

	// Using Index for its actual return value is fine
	idx := strings.Index(s, "hello")
	_ = idx

	// Comparison with values other than -1 and 0 is fine
	if strings.Index(s, "hello") > 0 {
		// checking if not at start
	}

	if strings.Index(s, "hello") == 0 {
		// checking if at start
	}
}
