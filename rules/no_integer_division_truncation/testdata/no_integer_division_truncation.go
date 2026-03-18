package fixtures

// Invalid: integer division that truncates to zero.
func intDivTruncatesSimple() {
	ratio := 1 / 2 // MATCH /integer division of 1 by 2 results in zero; use floating-point division if intended/
	_ = ratio
}

// Invalid: larger numerator still smaller than denominator.
func intDivTruncatesLarger() {
	x := 3 / 10 // MATCH /integer division of 3 by 10 results in zero; use floating-point division if intended/
	_ = x
}

// Invalid: hex literals.
func intDivTruncatesHex() {
	x := 0x1 / 0x10 // MATCH /integer division of 0x1 by 0x10 results in zero; use floating-point division if intended/
	_ = x
}

// Valid: numerator equals denominator (result is 1).
func intDivEqual() {
	x := 2 / 2
	_ = x
}

// Valid: numerator larger than denominator.
func intDivNoTruncation() {
	x := 10 / 3
	_ = x
}

// Valid: floating-point division.
func floatDiv() {
	ratio := 1.0 / 2.0
	_ = ratio
}

// Valid: division involving variables.
func varDiv(a int) {
	x := a / 2
	_ = x
}

// Valid: zero numerator (result is intentionally zero).
func zeroDivision() {
	x := 0 / 5
	_ = x
}

// Valid: negative literal is a unary expression, not a basic literal.
func negativeNumerator() {
	x := -1 / 2
	_ = x
}
