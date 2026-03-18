package fixtures

// Invalid: if returning true then returning false can be simplified.
func bad1(x int) bool {
	if x > 0 { // MATCH /redundant if/else returning bool; simplify to return the condition directly/
		return true
	}
	return false
}

// Invalid: if returning false then returning true can also be simplified.
func bad2(x int) bool {
	if x > 0 { // MATCH /redundant if/else returning bool; simplify to return the condition directly/
		return false
	}
	return true
}

// Invalid: with a more complex condition.
func bad3(x int, y int) bool {
	if x > 0 && y > 0 { // MATCH /redundant if/else returning bool; simplify to return the condition directly/
		return true
	}
	return false
}

// Valid: the function already returns the condition directly.
func good1(x int) bool {
	return x > 0
}

// Valid: the if body does not just return a bool literal.
func good2(x int) bool {
	if x > 0 {
		doSomething()
		return true
	}
	return false
}

// Valid: the next statement is not a return.
func good3(x int) bool {
	if x > 0 {
		return true
	}
	doSomething()
	return false
}

// Valid: both return the same value.
func good4(x int) bool {
	if x > 0 {
		return true
	}
	return true
}

// Valid: if with an else clause (not the pattern we are detecting).
func good5(x int) int {
	if x > 0 {
		return 1
	} else {
		return 2
	}
}

// Valid: if has init statement.
func good6() bool {
	if x := compute(); x > 0 {
		return true
	}
	return false
}

// Valid: returning non-bool values.
func good7(x int) int {
	if x > 0 {
		return 1
	}
	return 0
}

// helper stubs to make the file parse
func doSomething() {}
func compute() int { return 0 }
