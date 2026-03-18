package fixtures

// Invalid: short variable declaration followed by conditional reassignment.
func bad1(useDefault bool) string {
	x := "custom" // MATCH /merge conditional assignment into variable declaration of x/
	if useDefault {
		x = "default"
	}
	return x
}

// Invalid: var declaration with value followed by conditional reassignment.
func bad2(flag bool) int {
	var n = 10 // MATCH /merge conditional assignment into variable declaration of n/
	if flag {
		n = 20
	}
	return n
}

// Invalid: different types.
func bad3(useFallback bool) []string {
	result := []string{"a", "b"} // MATCH /merge conditional assignment into variable declaration of result/
	if useFallback {
		result = []string{"c"}
	}
	return result
}

// Valid: the if body has multiple statements.
func good1(useDefault bool) string {
	x := "custom"
	if useDefault {
		x = "default"
		doSomething()
	}
	return x
}

// Valid: the if has an else branch.
func good2(useDefault bool) string {
	x := "custom"
	if useDefault {
		x = "default"
	} else {
		x = "other"
	}
	return x
}

// Valid: the if body assigns a different variable.
func good3(useDefault bool) string {
	x := "custom"
	if useDefault {
		y := "default"
		_ = y
	}
	return x
}

// Valid: not followed by an if statement.
func good4() string {
	x := "custom"
	return x
}

// Valid: var declaration without initial value (handled by useMergedVarDecl).
func good5(flag bool) int {
	var n int
	if flag {
		n = 20
	}
	return n
}

// Valid: the if has an init statement.
func good6(flag bool) string {
	x := "custom"
	if y := compute(); y > 0 {
		x = "positive"
	}
	return x
}

// Valid: assignment uses += not =.
func good7(flag bool) int {
	x := 1
	if flag {
		x += 2
	}
	return x
}

// Valid: multiple LHS in the assignment.
func good8(flag bool) int {
	x := 1
	if flag {
		x, _ = twoReturns()
	}
	return x
}

// helper stubs
func doSomething()         {}
func compute() int         { return 0 }
func twoReturns() (int, error) { return 0, nil }
