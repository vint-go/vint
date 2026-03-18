package fixtures

// Invalid: modulo by 1 always returns 0
func isOddWrong(x int) bool {
	return x%1 != 0 // MATCH /x % 1 is always zero/
}

// Invalid: modulo by 1 in assignment
func moduloOneAssign(x int) int {
	r := x % 1 // MATCH /x % 1 is always zero/
	return r
}

// Invalid: modulo by 1 in if condition
func moduloOneCondition(x int) {
	if x%1 == 0 { // MATCH /x % 1 is always zero/
		_ = x
	}
}

// Valid: modulo by 2
func isOddCorrect(x int) bool {
	return x%2 != 0
}

// Valid: modulo by a variable
func moduloByVar(x, y int) int {
	return x % y
}

// Valid: modulo by a constant other than 1
func moduloByThree(x int) int {
	return x % 3
}

// Valid: modulo by 10
func lastDigit(x int) int {
	return x % 10
}
