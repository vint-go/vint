package fixtures

// Invalid: else block containing only a nested if can be simplified to else if.
// The skipBalanced default (true) means this should still trigger when the outer
// if body has more than one statement.

func bad1() {
	a := true
	b := true
	if a {
		doA()
		doA()
	} else { // MATCH /can replace else block with else if/
		if b {
			doB()
		}
	}
}

func bad2() {
	a := true
	b := true
	if a {
		doA()
		doA()
	} else { // MATCH /can replace else block with else if/
		if b {
			doB()
			doB()
		}
	}
}

func bad3() {
	a := true
	b := true
	if a {
		doA()
	} else { // MATCH /can replace else block with else if/
		if b {
			doB()
			doB()
		}
	}
}

// Valid: already using else if.
func good1() {
	a := true
	b := true
	if a {
		doA()
	} else if b {
		doB()
	}
}

// Valid: else block has more than just an if statement.
func good2() {
	a := true
	b := true
	if a {
		doA()
	} else {
		doC()
		if b {
			doB()
		}
	}
}

// Valid: balanced case (skipBalanced=true by default):
// both outer if and inner if have single statements.
func good3() {
	a := true
	b := true
	if a {
		doA()
	} else {
		if b {
			doB()
		}
	}
}

// Valid: no else branch.
func good4() {
	a := true
	if a {
		doA()
	}
}

// helper stubs to make the file parse
func doA() {}
func doB() {}
func doC() {}
