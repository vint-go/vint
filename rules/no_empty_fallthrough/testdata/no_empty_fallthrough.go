package fixtures

func handleOneOrTwo() {}
func handleDefault()  {}
func handleThree()    {}
func doSomething()    {}

func emptyFallthroughSimple() {
	x := 1
	switch x {
	case 1: // MATCH /case clause with only a fallthrough can be combined with the next case/
		fallthrough
	case 2:
		handleOneOrTwo()
	}
}

func emptyFallthroughToDefault() {
	x := 3
	switch x {
	case 3: // MATCH /case clause with only a fallthrough can be combined with the next case/
		fallthrough
	default:
		handleDefault()
	}
}

func validCombinedCases() {
	x := 1
	switch x {
	case 1, 2:
		handleOneOrTwo()
	}
}

func validDefaultOnly() {
	x := 1
	switch x {
	default:
		handleDefault()
	}
}

func validFallthroughWithLogic() {
	x := 1
	switch x {
	case 1:
		doSomething()
		fallthrough
	case 2:
		handleOneOrTwo()
	}
}

func validNormalCases() {
	x := 1
	switch x {
	case 1:
		handleOneOrTwo()
	case 2:
		handleThree()
	case 3:
		handleDefault()
	}
}
