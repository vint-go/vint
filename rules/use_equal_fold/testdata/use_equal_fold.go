package fixtures

import "strings"

func invalid(a, b, s string) {
	if strings.ToLower(a) == strings.ToLower(b) { // MATCH /use strings.EqualFold instead of strings.ToLower or strings.ToUpper for case-insensitive comparison/
		// case-insensitive comparison via ToLower
	}

	if strings.ToUpper(s) == "HELLO" { // MATCH /use strings.EqualFold instead of strings.ToLower or strings.ToUpper for case-insensitive comparison/
		// suboptimal case-insensitive comparison
	}

	if "hello" == strings.ToLower(s) { // MATCH /use strings.EqualFold instead of strings.ToLower or strings.ToUpper for case-insensitive comparison/
		// reversed comparison with ToLower
	}

	if strings.ToUpper(a) != strings.ToUpper(b) { // MATCH /use strings.EqualFold instead of strings.ToLower or strings.ToUpper for case-insensitive comparison/
		// inequality check with ToUpper
	}

	if strings.ToLower(s) != "world" { // MATCH /use strings.EqualFold instead of strings.ToLower or strings.ToUpper for case-insensitive comparison/
		// inequality check with ToLower
	}
}

func valid(a, b, s string) {
	if strings.EqualFold(a, b) {
		// efficient case-insensitive comparison
	}

	if strings.EqualFold(s, "HELLO") {
		// correct and efficient
	}

	if a == b {
		// normal string comparison, no case conversion
	}

	// Using ToLower for purposes other than comparison is fine
	lower := strings.ToLower(s)
	_ = lower

	upper := strings.ToUpper(s)
	_ = upper

	// ToLower/ToUpper in non-equality comparisons
	if strings.ToLower(s) > "abc" {
		// not an equality check
	}
}
