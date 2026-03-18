package fixtures

func sideEffect()      {}
func doSomething()     {}
func condition() bool  { return true }
func getValue() int    { return 0 }

var x int

// Invalid: non-assignment statement in if init clause
func invalidIfSideEffect() {
	if sideEffect(); condition() { // MATCH /avoid side effects in init clause of conditional statement; move the statement before the if or switch/
		doSomething()
	}
}

// Invalid: non-assignment statement in switch init clause
func invalidSwitchSideEffect() {
	switch sideEffect(); x { // MATCH /avoid side effects in init clause of conditional statement; move the statement before the if or switch/
	case 1:
		doSomething()
	}
}

// Valid: assignment in if init clause
func validIfAssignment() {
	if v := getValue(); v > 0 {
		doSomething()
	}
}

// Valid: assignment in switch init clause
func validSwitchAssignment() {
	switch v := getValue(); v {
	case 1:
		doSomething()
	}
}

// Valid: no init clause
func validNoInit() {
	if condition() {
		doSomething()
	}

	switch x {
	case 1:
		doSomething()
	}
}

// Valid: side effect moved before if
func validSideEffectBeforeIf() {
	sideEffect()
	if condition() {
		doSomething()
	}
}
