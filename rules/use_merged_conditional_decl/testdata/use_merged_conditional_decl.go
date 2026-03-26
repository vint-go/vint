package fixtures

// Invalid: short variable declaration with bool followed by opposite conditional reassignment.
func bad1(cond bool) bool {
	x := false // MATCH /merge conditional assignment into variable declaration of x/
	if cond {
		x = true
	}
	return x
}

// Invalid: var declaration with bool followed by opposite conditional reassignment.
func bad2(flag bool) bool {
	var n = true // MATCH /merge conditional assignment into variable declaration of n/
	if flag {
		n = false
	}
	return n
}

// Invalid: short variable declaration with true, reassigned to false.
func bad3(useFallback bool) bool {
	result := true // MATCH /merge conditional assignment into variable declaration of result/
	if useFallback {
		result = false
	}
	return result
}

// Valid: non-boolean types should not be flagged.
func good_string(useDefault bool) string {
	x := "custom"
	if useDefault {
		x = "default"
	}
	return x
}

// Valid: non-boolean int type should not be flagged.
func good_int(flag bool) int {
	var n = 10
	if flag {
		n = 20
	}
	return n
}

// Valid: non-boolean slice type should not be flagged.
func good_slice(useFallback bool) []string {
	result := []string{"a", "b"}
	if useFallback {
		result = []string{"c"}
	}
	return result
}

// Valid: the if body has multiple statements.
func good1(useDefault bool) bool {
	x := false
	if useDefault {
		x = true
		doSomething()
	}
	return x
}

// Valid: the if has an else branch.
func good2(useDefault bool) bool {
	x := false
	if useDefault {
		x = true
	} else {
		x = false
	}
	return x
}

// Valid: the if body assigns a different variable.
func good3(useDefault bool) bool {
	x := false
	if useDefault {
		y := true
		_ = y
	}
	return x
}

// Valid: not followed by an if statement.
func good4() bool {
	x := false
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
func good6(flag bool) bool {
	x := false
	if y := compute(); y > 0 {
		x = true
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

// Valid: both are booleans but reassignment is the SAME value, not opposite.
func good9(flag bool) bool {
	x := true
	if flag {
		x = true
	}
	return x
}

// Valid: initial value is boolean but reassignment is not a boolean literal.
func good10(flag bool) bool {
	x := false
	if flag {
		x = compute() > 0
	}
	return x
}

// helper stubs
func doSomething()              {}
func compute() int              { return 0 }
func twoReturns() (int, error) { return 0, nil }
