package fixtures

func handleOne()     {}
func handleTwo()     {}
func handleThree()   {}
func handleDefault() {}

func misplacedDefault(x int) {
	switch x {
	case 1:
		handleOne()
	default: // MATCH /default case should be the first or last case in the switch/
		handleDefault()
	case 2:
		handleTwo()
	}
}

func misplacedDefaultMiddleOf4(x int) {
	switch x {
	case 1:
		handleOne()
	case 2:
		handleTwo()
	default: // MATCH /default case should be the first or last case in the switch/
		handleDefault()
	case 3:
		handleThree()
	}
}

func defaultAtEnd(x int) {
	switch x {
	case 1:
		handleOne()
	case 2:
		handleTwo()
	default:
		handleDefault()
	}
}

func defaultAtBeginning(x int) {
	switch x {
	default:
		handleDefault()
	case 1:
		handleOne()
	case 2:
		handleTwo()
	}
}

func noDefaultCase(x int) {
	switch x {
	case 1:
		handleOne()
	case 2:
		handleTwo()
	}
}

func defaultOnlyTwo(x int) {
	switch x {
	default:
		handleDefault()
	case 1:
		handleOne()
	}
}

func defaultOnlySingle(x int) {
	switch x {
	default:
		handleDefault()
	}
}

type myInterface interface {
	method()
}

type myType struct{}

func (myType) method() {}

func typeSwitchMisplaced(v interface{}) {
	switch v.(type) {
	case int:
		handleOne()
	default: // MATCH /default case should be the first or last case in the switch/
		handleDefault()
	case string:
		handleTwo()
	}
}

func typeSwitchAtEnd(v interface{}) {
	switch v.(type) {
	case int:
		handleOne()
	case string:
		handleTwo()
	default:
		handleDefault()
	}
}

func typeSwitchAtBeginning(v interface{}) {
	switch v.(type) {
	default:
		handleDefault()
	case int:
		handleOne()
	case string:
		handleTwo()
	}
}
