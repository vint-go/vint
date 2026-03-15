package fixtures

// Invalid: Magic number 20 in multiplication operation.
func compute(y int) {
	_ = y * 20 // MATCH /magic number: 20, in <operation> detected/
}

// Invalid: Magic number 10 as left operand.
func compute2(y int) {
	_ = 10 * y // MATCH /magic number: 10, in <operation> detected/
}

// Invalid: Magic number 5 in a chained operation.
func compute3(y int) {
	_ = 5 * y // MATCH /magic number: 5, in <operation> detected/
}

// Invalid: Magic number 6 in a chained operation.
func compute3b(y int) {
	_ = y * 6 // MATCH /magic number: 6, in <operation> detected/
}

// Invalid: Magic number 42 in a parenthesized sub-expression.
func compute4a() {
	const c = 24
	_ = c + 42 // MATCH /magic number: 42, in <operation> detected/
}

// Invalid: Magic number 10 in a parenthesized sub-expression.
func compute4b() {
	const c = 24
	_ = c + 10 // MATCH /magic number: 10, in <operation> detected/
}

// Invalid: Magic number 42 in arithmetic operation inside a condition.
func compute5a(x int32) {
	if (42*x) > 0 { // MATCH /magic number: 42, in <operation> detected/
	}
}

// Invalid: Magic number 10 as comparison operand in condition.
func compute5b(x int32) {
	if x > 10 { // MATCH /magic number: 10, in <operation> detected/
	}
}

// Invalid: Magic number 10 as comparison operand in condition (1.0 excluded).
func compute6(x float32) {
	if 10 < x { // MATCH /magic number: 10, in <operation> detected/
	}
}

// Valid: Using a named constant in an operation.
const multiplier = 20

func computeGood(y int) {
	_ = y * multiplier
}

// Valid: Using named constants in operations.
const (
	factor = 42
	base   = 10
)

func compute2Good() {
	const c = 24
	_ = c + (factor * base)
}

// Valid: 1 and 0 are excluded by default.
func computeSimple(y int) {
	_ = y * 1
	_ = y + 0
}

// Valid: 1.0 is excluded by default.
func computeFloat(x float32) {
	if x > 1.0 {
	}
}
