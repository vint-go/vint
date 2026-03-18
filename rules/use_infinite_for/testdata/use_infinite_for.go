package fixtures

// Invalid: for true { ... } should be for { ... }
func badForTrue() {
	for true { // MATCH /use for { ... } instead of for true { ... }/
		break
	}
}

// Invalid: for true with multiple statements
func badForTrueMultiple() {
	x := 0
	for true { // MATCH /use for { ... } instead of for true { ... }/
		x++
		if x > 10 {
			break
		}
	}
}

// Valid: for { ... } (already idiomatic)
func goodInfiniteFor() {
	for {
		break
	}
}

// Valid: for with a non-true condition
func goodForWithCondition() {
	x := 0
	for x < 10 {
		x++
	}
}

// Valid: standard 3-clause for loop
func goodStandardFor() {
	for i := 0; i < 10; i++ {
	}
}

// Valid: for range loop
func goodRangeFor() {
	items := []int{1, 2, 3}
	for range items {
	}
}

// Valid: for with init and true condition (not a simple for true)
func goodForInitTrue() {
	for i := 0; true; i++ {
		if i > 10 {
			break
		}
	}
}

// Valid: for false (not true)
func goodForFalse() {
	for false {
		break
	}
}
