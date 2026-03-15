package fixtures

func doubleNegation(ready bool) {
	if !(!(ready)) { // MATCH /can simplify !(!(ready)) to ready/
		process()
	}
}

func negatedComparison(x, y int) {
	if !(x >= y) { // MATCH /can simplify !(x >= y) to x < y/
		handle()
	}
	if !(x < y) { // MATCH /can simplify !(x < y) to x >= y/
		handle()
	}
	if !(x == y) { // MATCH /can simplify !(x == y) to x != y/
		handle()
	}
	if !(x != y) { // MATCH /can simplify !(x != y) to x == y/
		handle()
	}
	if !(x > y) { // MATCH /can simplify !(x > y) to x <= y/
		handle()
	}
	if !(x <= y) { // MATCH /can simplify !(x <= y) to x > y/
		handle()
	}
}

func combinedComparisons(x int) {
	if x > 10 || x == 10 { // MATCH /can simplify x > 10 || x == 10 to x >= 10/
		handle()
	}
	if x < 5 || x == 5 { // MATCH /can simplify x < 5 || x == 5 to x <= 5/
		handle()
	}
	if x == 10 || x > 10 { // MATCH /can simplify x > 10 || x == 10 to x >= 10/
		handle()
	}
}

func incrementDecrement(x, y int) {
	if x > y-1 { // MATCH /can simplify x > y-1 to x >= y/
		handle()
	}
	if x < y+1 { // MATCH /can simplify x < y+1 to x <= y/
		handle()
	}
	if x >= y+1 { // MATCH /can simplify x >= y+1 to x > y/
		handle()
	}
	if x <= y-1 { // MATCH /can simplify x <= y-1 to x < y/
		handle()
	}
}

func rangeFolding(x int) {
	if x >= 10 && x <= 10 { // MATCH /can simplify x >= 10 && x <= 10 to x == 10/
		handle()
	}
	if x <= 5 && x >= 5 { // MATCH /can simplify x >= 5 && x <= 5 to x == 5/
		handle()
	}
}

func negatedEquality(a, b bool) {
	if !(a) == !(b) { // MATCH /can simplify !((a)) == !((b)) to ((a)) == ((b))/
		handle()
	}
}

// Valid cases - no matches expected

func validSimpleCondition(ready bool) {
	if ready {
		process()
	}
}

func validSimpleComparison(x, y int) {
	if x < y {
		handle()
	}
}

func validCombinedComparison(x int) {
	if x >= 10 {
		handle()
	}
}

func validDifferentOperands(x, y int) {
	// Different operands on each side, not a simplifiable pattern
	if x > 10 || y == 10 {
		handle()
	}
}

func validDifferentValues(x int) {
	// Different values, not simplifiable
	if x > 10 || x == 5 {
		handle()
	}
}

func validSingleNegation(ready bool) {
	if !ready {
		handle()
	}
}

func validFloatAvoidance(x float64) {
	// Floats should be avoided
	if !(x >= 1.0) {
		handle()
	}
}

func process() {}
func handle()  {}
