package fixtures

import "strings"

func invalidPrefix(s string) string {
	if strings.HasPrefix(s, "http://") { // MATCH /replace HasPrefix and manual slicing with strings.TrimPrefix/
		s = s[len("http://"):]
	}
	return s
}

func invalidSuffix(s string) string {
	if strings.HasSuffix(s, ".go") { // MATCH /replace HasSuffix and manual slicing with strings.TrimSuffix/
		s = s[:len(s)-len(".go")]
	}
	return s
}

func invalidPrefixVariable(s, prefix string) string {
	if strings.HasPrefix(s, prefix) { // MATCH /replace HasPrefix and manual slicing with strings.TrimPrefix/
		s = s[len(prefix):]
	}
	return s
}

func invalidSuffixVariable(s, suffix string) string {
	if strings.HasSuffix(s, suffix) { // MATCH /replace HasSuffix and manual slicing with strings.TrimSuffix/
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func validTrimPrefix(s string) string {
	return strings.TrimPrefix(s, "http://")
}

func validTrimSuffix(s string) string {
	return strings.TrimSuffix(s, ".go")
}

func validMultipleStatements(s string) string {
	// Body has more than one statement - not a simple trim
	if strings.HasPrefix(s, "http://") {
		s = s[len("http://"):]
		s = "https://" + s
	}
	return s
}

func validWithElse(s string) string {
	// Has else branch - not a simple trim
	if strings.HasPrefix(s, "http://") {
		s = s[len("http://"):]
	} else {
		s = "default"
	}
	return s
}

func validDifferentVariable(s, other string) string {
	// Assigns to a different variable than the one checked
	if strings.HasPrefix(s, "http://") {
		other = s[len("http://"):]
	}
	_ = other
	return s
}

func validDifferentSliceTarget(s, other string) string {
	// Slices a different variable than the one checked
	if strings.HasPrefix(s, "http://") {
		s = other[len("http://"):]
	}
	return s
}

func validNotSliceExpr(s string) string {
	// RHS is not a slice expression
	if strings.HasPrefix(s, "http://") {
		s = "replaced"
	}
	return s
}
