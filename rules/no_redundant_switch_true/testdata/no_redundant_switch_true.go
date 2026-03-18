package fixtures

func handlePositive() {}
func handleNegative() {}
func handleZero()     {}
func doSomething()    {}

// Invalid: switch true is redundant
func redundantSwitchTrue() {
	x := 1
	switch true { // MATCH /redundant switch true; use switch without a tag instead/
	case x > 0:
		handlePositive()
	case x < 0:
		handleNegative()
	}
}

// Invalid: switch true with default
func redundantSwitchTrueWithDefault() {
	x := 1
	switch true { // MATCH /redundant switch true; use switch without a tag instead/
	case x > 0:
		handlePositive()
	default:
		handleNegative()
	}
}

// Valid: switch without tag (tagless switch)
func taglessSwitch() {
	x := 1
	switch {
	case x > 0:
		handlePositive()
	case x < 0:
		handleNegative()
	}
}

// Valid: switch with a non-true expression tag
func switchWithExprTag() {
	x := 1
	switch x {
	case 1:
		handlePositive()
	case -1:
		handleNegative()
	}
}

// Valid: switch with a variable named true (not the builtin)
func switchWithBoolVar() {
	x := true
	switch x {
	case true:
		handlePositive()
	case false:
		handleNegative()
	}
}
