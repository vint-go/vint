package fixtures

// Invalid: len(s) >= 0 is always true.
func alwaysTrueGEQ(s []int) {
	if len(s) >= 0 { // MATCH /len(x) >= 0 is always true/
		doSomething()
	}
}

// Invalid: len(s) < 0 is always false.
func alwaysFalseLSS(s []int) {
	if len(s) < 0 { // MATCH /len(x) < 0 is always false/
		doSomething()
	}
}

// Invalid: 0 <= len(s) is equivalent to len(s) >= 0, always true.
func alwaysTrueFlipped(s string) {
	if 0 <= len(s) { // MATCH /len(x) >= 0 is always true/
		doSomething()
	}
}

// Invalid: 0 > len(s) is equivalent to len(s) < 0, always false.
func alwaysFalseFlipped(s string) {
	if 0 > len(s) { // MATCH /len(x) < 0 is always false/
		doSomething()
	}
}

// Valid: len(s) > 0 is a meaningful check.
func validGreaterThanZero(s []int) {
	if len(s) > 0 {
		doSomething()
	}
}

// Valid: len(s) == 0 is a meaningful check.
func validEqualZero(s []int) {
	if len(s) == 0 {
		doSomething()
	}
}

// Valid: len(s) != 0 is a meaningful check.
func validNotEqualZero(s []int) {
	if len(s) != 0 {
		doSomething()
	}
}

// Valid: len(s) <= 0 is not always-true/false (could be == 0).
func validLEQZero(s []int) {
	if len(s) <= 0 {
		doSomething()
	}
}

// Valid: comparing len against a non-zero value.
func validNonZero(s []int) {
	if len(s) >= 1 {
		doSomething()
	}
}

func doSomething() {}
