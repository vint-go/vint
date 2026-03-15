package fixtures

// Invalid: inverted loop condition.
func invertedLoop(n int) {
	for i := 0; i > n; i++ { // MATCH /suspicious loop condition: loop variable increments but condition uses >/
		_ = i
	}
}

// Invalid: impossible range condition.
func impossibleRange(x int) {
	if x < 5 && x > 10 { // MATCH /impossible condition: x cannot satisfy both comparisons simultaneously/
		_ = x
	}
}

// Invalid: x cannot equal two different constants.
func doubleEquality(x int) {
	if x == 1 && x == 2 { // MATCH /impossible condition: x cannot equal both values simultaneously/
		_ = x
	}
}

// Invalid: inverted loop with >=.
func invertedLoopGEQ(n int) {
	for i := 0; i >= n; i++ { // MATCH /suspicious loop condition: loop variable increments but condition uses >=/
		_ = i
	}
}

// Invalid: impossible range with <=.
func impossibleRangeLEQ(x int) {
	if x <= 3 && x > 10 { // MATCH /impossible condition: x cannot satisfy both comparisons simultaneously/
		_ = x
	}
}

// Invalid: impossible range, both strict.
func impossibleRangeStrict(x int) {
	if x < 5 && x > 5 { // MATCH /impossible condition: x cannot satisfy both comparisons simultaneously/
		_ = x
	}
}

// Valid: correct loop condition.
func validLoop(n int) {
	for i := 0; i < n; i++ {
		_ = i
	}
}

// Valid: valid range check.
func validRange(x int) {
	if x > 5 && x < 10 {
		_ = x
	}
}

// Valid: same equality value.
func sameEquality(x int) {
	if x == 1 && x == 1 {
		_ = x
	}
}

// Valid: different variables in equality.
func differentVars(x, y int) {
	if x == 1 && y == 2 {
		_ = x
	}
}

// Valid: OR condition (not AND).
func orCondition(x int) {
	if x < 5 || x > 10 {
		_ = x
	}
}

// Valid: loop with !=.
func validLoopNEQ(n int) {
	for i := 0; i != n; i++ {
		_ = i
	}
}

// Valid: non-contradictory range.
func validRangeLEQGEQ(x int) {
	if x <= 10 && x >= 5 {
		_ = x
	}
}
