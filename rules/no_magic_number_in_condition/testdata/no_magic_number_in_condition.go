package fixtures

// Invalid: Magic number 7 in if condition.
func checkValue(x int) {
	if x > 7 { // MATCH /magic number: 7, in <condition> detected/
		// do something
	}
}

// Invalid: Magic number 8 in if condition (left operand).
func checkValue2(x int) {
	if 8 > x { // MATCH /magic number: 8, in <condition> detected/
		// do something
	}
}

// Invalid: Magic number 42 in equality check.
func checkEquality(x int) {
	if x == 42 { // MATCH /magic number: 42, in <condition> detected/
		// do something
	}
}

// Invalid: Magic number 3.14 in float condition.
func checkFloat(x float64) {
	if x > 3.14 { // MATCH /magic number: 3.14, in <condition> detected/
		// do something
	}
}

// Valid: Using a named constant in condition.
const maxRetries = 7

func checkWithConst(x int) {
	if x > maxRetries {
		// do something
	}
}

// Valid: 1.0 is excluded by default.
func checkDefaultExcluded(x float32) {
	if x > 1.0 {
		// do something
	}
}

// Valid: 0.0 is excluded by default.
func checkZeroExcluded(x float32) {
	if x < 0.0 {
		// do something
	}
}

// Valid: 0 is excluded by default.
func checkZeroInt(x int) {
	if x == 0 {
		// do something
	}
}

// Valid: 1 is excluded by default.
func checkOneInt(x int) {
	if x == 1 {
		// do something
	}
}

// Valid: String comparisons are not checked.
func checkString(y string) {
	if "test" == y {
		// do something
	}
}

// Valid: Boolean conditions are not checked.
func checkBool(b bool) {
	if b {
		// do something
	}
}
