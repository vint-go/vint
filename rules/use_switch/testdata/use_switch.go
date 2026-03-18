package fixtures

func handleOne()   {}
func handleTwo()   {}
func handleThree() {}
func handleFour()  {}
func handleDefault() {}

func invalidIfElseChain(x int) {
	if x == 1 { // MATCH /could replace if-else chain with 3 branches by a switch statement/
		handleOne()
	} else if x == 2 {
		handleTwo()
	} else if x == 3 {
		handleThree()
	}
}

func invalidIfElseChainWithElse(x int) {
	if x == 1 { // MATCH /could replace if-else chain with 3 branches by a switch statement/
		handleOne()
	} else if x == 2 {
		handleTwo()
	} else if x == 3 {
		handleThree()
	} else {
		handleDefault()
	}
}

func invalidIfElseChainFour(x int) {
	if x == 1 { // MATCH /could replace if-else chain with 4 branches by a switch statement/
		handleOne()
	} else if x == 2 {
		handleTwo()
	} else if x == 3 {
		handleThree()
	} else if x == 4 {
		handleFour()
	}
}

func invalidTwoBranches(x int) {
	if x == 1 { // MATCH /could replace if-else chain with 2 branches by a switch statement/
		handleOne()
	} else if x == 2 {
		handleTwo()
	}
}

// Valid: single if without else-if
func validSingleIf(x int) {
	if x == 1 {
		handleOne()
	}
}

// Valid: if-else without else-if
func validIfElse(x int) {
	if x == 1 {
		handleOne()
	} else {
		handleDefault()
	}
}

// Valid: already uses a switch
func validSwitch(x int) {
	switch x {
	case 1:
		handleOne()
	case 2:
		handleTwo()
	case 3:
		handleThree()
	}
}
