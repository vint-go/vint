package fixtures

func doSomething()     {}
func doSomethingElse() {}
func handleOne()       {}
func handleTwo()       {}
func use(x int)        {}
func computeValue() int { return 1 }

// Invalid: unnecessary block in function body
func unnecessaryBlockInFunc() {
	{ // MATCH /unnecessary block statement can be removed/
		doSomething()
		doSomethingElse()
	}
}

// Invalid: unnecessary block in switch case
func unnecessaryBlockInCase() {
	x := 1
	switch x {
	case 1:
		{ // MATCH /unnecessary block in case clause/
			handleOne()
		}
	}
}

// Invalid: empty block statement (no declarations, no short var decls)
func unnecessaryEmptyBlock() {
	{ // MATCH /unnecessary block statement can be removed/
	}
}

// Invalid: block with only plain assignments (not short var decl)
func unnecessaryBlockPlainAssign() {
	x := 0
	{ // MATCH /unnecessary block statement can be removed/
		x = 1
		use(x)
	}
}

// Valid: block is needed for scoping with short variable declaration
func validBlockForScoping() {
	{
		x := computeValue()
		use(x)
	}
}

// Valid: block with var declaration
func validBlockWithVarDecl() {
	{
		var x int
		x = computeValue()
		use(x)
	}
}

// Valid: normal function body without nested blocks
func validNoNestedBlock() {
	doSomething()
	doSomethingElse()
}

// Valid: switch case without block
func validCaseNoBlock() {
	x := 1
	switch x {
	case 1:
		handleOne()
	case 2:
		handleTwo()
	}
}
