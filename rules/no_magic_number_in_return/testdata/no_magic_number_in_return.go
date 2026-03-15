package fixtures

// Invalid: Magic number 3 in return statement.
func getCount() int {
	return 3 // MATCH /magic number: 3, in <return> detected/
}

// Invalid: Magic number 2.0 in return statement.
func getPi() float64 {
	return 2.0 // MATCH /magic number: 2, in <return> detected/
}

// Invalid: Magic number 42 in binary expression in return statement.
func addOffset(x int32) int32 {
	return x + 42 // MATCH /magic number: 42, in <return> detected/
}

// Invalid: Magic number 42 on left side of binary expression in return.
func addOffset2(x int32) int32 {
	return 42 + x // MATCH /magic number: 42, in <return> detected/
}

// Invalid: Magic number 42 in parenthesized expression in return.
func complexReturn(x int32) int32 {
	return x + (42 * 1) // MATCH /magic number: 42, in <return> detected/
}

// Invalid: Magic number 10 in parenthesized expression in return.
func complexReturn2(x int32) int32 {
	return (1 * x) + 10 // MATCH /magic number: 10, in <return> detected/
}

// Invalid: Magic number 3.0 (float) in return statement.
func getFloat() float32 {
	return 3.0 // MATCH /magic number: 3, in <return> detected/
}

// Valid: Using a named constant in return.
const defaultCount = 3

func getCountGood() int {
	return defaultCount
}

// Valid: Using a named constant in binary expression in return.
const offset = 42

func addOffsetGood(x int32) int32 {
	return x + offset
}

// Valid: String literals are not checked.
func getString() string {
	return "3"
}

// Valid: 0.0 is excluded by default.
func getZero() float32 {
	return 0.0
}

// Valid: 1.0 is excluded by default.
func getOne() float32 {
	return 1.0
}

// Valid: 0 is excluded by default.
func getZeroInt() int {
	return 0
}

// Valid: 1 is excluded by default.
func getOneInt() int {
	return 1
}
