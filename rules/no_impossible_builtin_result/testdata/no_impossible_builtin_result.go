package fixtures

// Invalid: len() never returns negative, so len(s) < -1 is always false.
func lenLessThanNeg(s []int) {
	if len(s) < -1 { // MATCH /len() never returns a negative value, comparison always evaluates to false/
		doSomething()
	}
}

// Invalid: cap() never returns negative, so cap(s) < -1 is always false.
func capLessThanNeg(s []int) {
	if cap(s) < -1 { // MATCH /cap() never returns a negative value, comparison always evaluates to false/
		doSomething()
	}
}

// Invalid: len() never returns negative, so len(s) <= -1 is always false.
func lenLEQNeg(s []int) {
	if len(s) <= -1 { // MATCH /len() never returns a negative value, comparison always evaluates to false/
		doSomething()
	}
}

// Invalid: len() never returns negative, so len(s) == -1 is always false.
func lenEqualNeg(s []int) {
	if len(s) == -1 { // MATCH /len() never returns a negative value, comparison always evaluates to false/
		doSomething()
	}
}

// Invalid: len() never returns negative, so len(s) > -1 is always true.
func lenGreaterThanNeg(s []int) {
	if len(s) > -1 { // MATCH /len() never returns a negative value, comparison always evaluates to true/
		doSomething()
	}
}

// Invalid: len() never returns negative, so len(s) >= -1 is always true.
func lenGEQNeg(s []int) {
	if len(s) >= -1 { // MATCH /len() never returns a negative value, comparison always evaluates to true/
		doSomething()
	}
}

// Invalid: len() never returns negative, so len(s) != -1 is always true.
func lenNEQNeg(s []int) {
	if len(s) != -1 { // MATCH /len() never returns a negative value, comparison always evaluates to true/
		doSomething()
	}
}

// Invalid: flipped comparison, -1 > len(s) is equivalent to len(s) < -1.
func flippedNegLenComparison(s []int) {
	if -1 > len(s) { // MATCH /len() never returns a negative value, comparison always evaluates to false/
		doSomething()
	}
}

// Invalid: cap() with larger negative value.
func capLargerNeg(s []int) {
	if cap(s) <= -100 { // MATCH /cap() never returns a negative value, comparison always evaluates to false/
		doSomething()
	}
}

// Valid: len(s) == 0 is a meaningful check.
func validLenEqualZero(s []int) {
	if len(s) == 0 {
		doSomething()
	}
}

// Valid: len(s) > 0 is a meaningful check.
func validLenGreaterZero(s []int) {
	if len(s) > 0 {
		doSomething()
	}
}

// Valid: cap(s) > 0 is a meaningful check.
func validCapGreaterZero(s []int) {
	if cap(s) > 0 {
		doSomething()
	}
}

// Valid: len(s) < 5 is a meaningful check.
func validLenLessThanPositive(s []int) {
	if len(s) < 5 {
		doSomething()
	}
}

// Valid: comparing against a positive number.
func validCapEqualPositive(s []int) {
	if cap(s) == 10 {
		doSomething()
	}
}

func doSomething() {}
