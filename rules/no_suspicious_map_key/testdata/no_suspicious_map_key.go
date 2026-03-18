package fixtures

func getValue() string {
	return "key"
}

func suspiciousMapKeys() {
	// Invalid: trailing space in key
	_ = map[string]int{
		"foo":  1,
		"bar ": 2, // MATCH /suspicious map key "bar " has trailing space/
	}

	// Invalid: leading space in key
	_ = map[string]int{
		" foo": 1, // MATCH /suspicious map key " foo" has leading space/
		"bar":  2,
	}

	// Invalid: duplicate non-literal key
	_ = map[string]int{
		getValue(): 1,
		getValue(): 2, // MATCH /suspicious duplicate key getValue() in map literal/
	}

	// Valid: normal keys without space issues
	_ = map[string]int{
		"foo": 1,
		"bar": 2,
	}

	// Valid: intentional multi-space padding (not single space)
	_ = map[string]int{
		"  foo": 1,
		"bar  ": 2,
	}

	// Valid: integer key map, no whitespace checks
	_ = map[int]string{
		1: "a",
		2: "b",
	}

	// Valid: different non-literal keys
	_ = map[string]int{
		getValue(): 1,
	}
}
