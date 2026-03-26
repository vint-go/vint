package fixtures

const ModeDebug = "debug"

// Invalid: "debug" appears 2 times, meeting custom min-occurrences of 2.

func IsDebug(mode string) bool {
	return mode == "debug" // MATCH /string literal "debug" matches constant ModeDebug, use the constant instead/
}

func IsDebug2(mode string) bool {
	return mode == "debug" // MATCH /string literal "debug" matches constant ModeDebug, use the constant instead/
}

// Valid: "debug" used as a map key should not be counted.

func MapKeyDebug() map[string]string {
	return map[string]string{
		"debug": "true",
	}
}
