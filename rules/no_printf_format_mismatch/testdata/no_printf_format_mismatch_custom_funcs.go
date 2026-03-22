package fixtures

import "os"

// This file tests custom printf function matching by configuring
// os.Setenv as a custom printf-like function (for testing purposes only).

// Invalid: format string expects 2 args, but only 1 provided.
func customFuncArgCountMismatch() {
	os.Setenv("count: %d name: %s", "only_one") // MATCH /format string expects 2 argument(s), but 1 provided/
}

// Invalid: format verb %d expects integer, got string.
func customFuncVerbTypeMismatch() {
	os.Setenv("count: %d", "not_an_int") // MATCH /format verb %d expects an integer type, got string (argument #1)/
}

// Valid: format string matches the provided argument.
func customFuncValid() {
	os.Setenv("name: %s", "Alice")
}

// Valid: no format verbs and no extra args beyond the format string.
func customFuncNoVerbs() {
	os.Setenv("plain string")
}
