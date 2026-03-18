package fixtures

// Invalid: name parameter is overwritten before use.
func process(name string) string {
	name = "default" // MATCH /argument 'name' is overwritten before first use/
	return name
}

// Invalid: multiple parameters, one overwritten.
func multiParam(a int, b int) int {
	a = 10 // MATCH /argument 'a' is overwritten before first use/
	return a + b
}

// Valid: parameter is read before being overwritten.
func readFirst(name string) string {
	if name == "" {
		name = "default"
	}
	return name
}

// Valid: parameter is read in the first statement.
func returnParam(x int) int {
	return x
}

// Valid: parameter is used in a call.
func callWithParam(x int) {
	println(x)
}

// Valid: blank identifier parameter.
func blankParam(_ int) {
}

// Valid: short var decl shadows the param, not an overwrite.
func shortVarDecl(x int) int {
	x := 5
	return x
}

// Valid: parameter is used in the RHS of the assignment.
func usedInRHS(x int) int {
	x = x + 1
	return x
}

// Valid: parameter is read in an if condition.
func readInIf(name string) string {
	if name != "" {
		return name
	}
	name = "default"
	return name
}

// Valid: parameter read in a for loop.
func readInFor(n int) int {
	sum := 0
	for i := 0; i < n; i++ {
		sum += i
	}
	return sum
}

// Valid: parameter is used in a range expression.
func readInRange(items []int) int {
	sum := 0
	for _, v := range items {
		sum += v
	}
	return sum
}

// Valid: parameter used in a switch tag.
func readInSwitch(x int) string {
	switch x {
	case 1:
		return "one"
	default:
		return "other"
	}
}

// Valid: parameter read in a goroutine.
func readInGo(x int) {
	go func() {
		_ = x
	}()
}

// Valid: parameter read in a defer.
func readInDefer(x int) {
	defer println(x)
}

// Valid: no parameters.
func noParams() int {
	return 42
}
