package fixtures

func handleOne()      {}
func handleTwo()      {}
func handleThree()    {}
func handleOneAgain() {}

func duplicateSwitchCase() {
	x := 1
	switch x {
	case 1:
		handleOne()
	case 2:
		handleTwo()
	case 1: // MATCH /duplicate case 1 in switch statement/
		handleOneAgain()
	}
}

func duplicateSwitchCaseString() {
	s := "a"
	switch s {
	case "a":
		handleOne()
	case "b":
		handleTwo()
	case "a": // MATCH /duplicate case "a" in switch statement/
		handleOneAgain()
	}
}

func validSwitchCase() {
	x := 1
	switch x {
	case 1:
		handleOne()
	case 2:
		handleTwo()
	case 3:
		handleThree()
	}
}

func duplicateSelectCase() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	select {
	case <-ch1:
		handleOne()
	case <-ch2:
		handleTwo()
	case <-ch1: // MATCH /duplicate case <-ch1 in select statement/
		handleOneAgain()
	}
}

func validSelectCase() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	select {
	case <-ch1:
		handleOne()
	case <-ch2:
		handleTwo()
	}
}

func duplicateTypeSwitchCase() {
	var x interface{} = 1
	switch x.(type) {
	case int:
		handleOne()
	case string:
		handleTwo()
	case int: // MATCH /duplicate case int in switch statement/
		handleOneAgain()
	}
}

func validTypeSwitchCase() {
	var x interface{} = 1
	switch x.(type) {
	case int:
		handleOne()
	case string:
		handleTwo()
	case float64:
		handleThree()
	}
}

func validSwitchCaseDefault() {
	x := 1
	switch x {
	case 1:
		handleOne()
	case 2:
		handleTwo()
	default:
		handleThree()
	}
}

func validSelectCaseDefault() {
	ch1 := make(chan int)
	select {
	case <-ch1:
		handleOne()
	default:
		handleTwo()
	}
}
