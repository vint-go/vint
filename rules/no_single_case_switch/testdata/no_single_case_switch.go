package fixtures

func handleOne()   {}
func handleTwo()   {}
func doSomething() {}

func singleCaseSwitch() {
	x := 1
	switch x { // MATCH /switch with a single case can be rewritten as an if statement/
	case 1:
		handleOne()
	}
}

func defaultOnlySwitch() {
	switch { // MATCH /switch with only a default case is redundant/
	default:
		doSomething()
	}
}

func defaultOnlySwitchWithTag() {
	x := 1
	switch x { // MATCH /switch with only a default case is redundant/
	default:
		doSomething()
	}
}

// Valid: multiple cases
func multiCaseSwitch() {
	x := 1
	switch x {
	case 1:
		handleOne()
	case 2:
		handleTwo()
	}
}

// Valid: single case with break
func singleCaseWithBreak() {
	x := 1
	switch x {
	case 1:
		handleOne()
		break
	}
}

// Valid: default-only with break
func defaultOnlyWithBreak() {
	switch {
	default:
		doSomething()
		break
	}
}

// Valid: if statement
func ifStatement() {
	x := 1
	if x == 1 {
		handleOne()
	}
}

// Valid: plain statement
func plainStatement() {
	doSomething()
}

// Valid: case plus default
func caseAndDefault() {
	x := 1
	switch x {
	case 1:
		handleOne()
	default:
		doSomething()
	}
}
