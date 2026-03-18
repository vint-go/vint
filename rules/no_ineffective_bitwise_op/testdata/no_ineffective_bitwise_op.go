package fixtures

// Invalid: XOR with 0 (x ^ 0 is always x)
func xorWithZero(x int) int {
	return x ^ 0 // MATCH /ineffective bitwise XOR with 0: the expression always equals the other operand/
}

// Invalid: XOR with 0 on the left (0 ^ x is always x)
func xorWithZeroLeft(x int) int {
	return 0 ^ x // MATCH /ineffective bitwise XOR with 0: the expression always equals the other operand/
}

// Invalid: OR with 0 (x | 0 is always x)
func orWithZero(x int) int {
	return x | 0 // MATCH /ineffective bitwise OR with 0: the expression always equals the other operand/
}

// Invalid: OR with 0 on the left (0 | x is always x)
func orWithZeroLeft(x int) int {
	return 0 | x // MATCH /ineffective bitwise OR with 0: the expression always equals the other operand/
}

// Invalid: AND with 0 (x & 0 is always 0)
func andWithZero(x int) int {
	return x & 0 // MATCH /ineffective bitwise AND with 0: the expression always equals 0/
}

// Invalid: AND with 0 on the left (0 & x is always 0)
func andWithZeroLeft(x int) int {
	return 0 & x // MATCH /ineffective bitwise AND with 0: the expression always equals 0/
}

// Invalid: AND NOT with 0 on the right (x &^ 0 is always x)
func andNotWithZero(x int) int {
	return x &^ 0 // MATCH /ineffective bitwise AND NOT with 0: the expression always equals the other operand/
}

// Invalid: AND NOT with 0 on the left (0 &^ x is always 0)
func andNotWithZeroLeft(x int) int {
	return 0 &^ x // MATCH /ineffective bitwise AND NOT with 0: the expression always equals 0/
}

// Valid: meaningful bitwise XOR
func validXor(x int) int {
	return x ^ 0xFF
}

// Valid: meaningful bitwise OR
func validOr(x int) int {
	return x | 0x0F
}

// Valid: meaningful bitwise AND
func validAnd(x int) int {
	return x & 0xFF
}

// Valid: meaningful bitwise AND NOT
func validAndNot(x int) int {
	return x &^ 0x0F
}

// Valid: non-zero literal on both sides
func validNonZero(x int) int {
	return 3 ^ 5
}
