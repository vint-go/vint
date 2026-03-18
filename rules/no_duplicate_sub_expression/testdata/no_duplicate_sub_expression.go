package fixtures

func doSomething() {}

type item struct {
	v int
}

// Invalid: same subexpression on both sides of comparison
func duplicateComparison(xs []item, i int) {
	if xs[i].v < xs[i].v { // MATCH /suspicious identical sub-expressions on both sides of '<'/
		doSomething()
	}
}

// Invalid: identical operands in logical expression
func duplicateLogicalOr(a bool) {
	if a || a { // MATCH /suspicious identical sub-expressions on both sides of '||'/
		doSomething()
	}
}

// Invalid: identical operands in logical AND
func duplicateLogicalAnd(a bool) {
	if a && a { // MATCH /suspicious identical sub-expressions on both sides of '&&'/
		doSomething()
	}
}

// Invalid: identical operands in equality
func duplicateEquality(x int) {
	if x == x { // MATCH /suspicious identical sub-expressions on both sides of '=='/
		doSomething()
	}
}

// Invalid: identical operands in inequality
func duplicateInequality(x int) {
	if x != x { // MATCH /suspicious identical sub-expressions on both sides of '!='/
		doSomething()
	}
}

// Invalid: identical operands in subtraction
func duplicateSubtraction(x int) {
	_ = x - x // MATCH /suspicious identical sub-expressions on both sides of '-'/
}

// Invalid: identical operands in division
func duplicateDivision(x int) {
	_ = x / x // MATCH /suspicious identical sub-expressions on both sides of '/'/
}

// Invalid: identical operands in modulo
func duplicateModulo(x int) {
	_ = x % x // MATCH /suspicious identical sub-expressions on both sides of '%'/
}

// Invalid: identical operands in bitwise AND
func duplicateBitwiseAnd(x int) {
	_ = x & x // MATCH /suspicious identical sub-expressions on both sides of '&'/
}

// Invalid: identical operands in bitwise OR
func duplicateBitwiseOr(x int) {
	_ = x | x // MATCH /suspicious identical sub-expressions on both sides of '|'/
}

// Invalid: identical operands in bitwise XOR
func duplicateBitwiseXor(x int) {
	_ = x ^ x // MATCH /suspicious identical sub-expressions on both sides of '^'/
}

// Valid: different indices
func validDifferentIndex(xs []item, i, j int) {
	if xs[i].v < xs[j].v {
		doSomething()
	}
}

// Valid: different operands in logical expression
func validLogicalOr(a, b bool) {
	if a || b {
		doSomething()
	}
}

// Valid: addition of same value (intentional doubling)
func validAddition(x int) {
	_ = x + x
}

// Valid: multiplication of same value (intentional squaring)
func validMultiplication(x int) {
	_ = x * x
}

// Valid: different operands in comparison
func validComparison(x, y int) {
	if x < y {
		doSomething()
	}
}

// Valid: shift with same value (intentional)
func validShift(x int) {
	_ = x << x
}
