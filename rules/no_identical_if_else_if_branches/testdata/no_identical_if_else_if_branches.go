package fixtures

func identicalIfElseIfBranches() {

	if true { // MATCH /"if...else if" chain with identical branches (lines 5 and 13)/
		print("something")
	} else if true {
		print("something else")
	} else if true {
		print("other thing")
	} else if false {
		println()
	} else {
		print("something")
	}

	if true { // MATCH /"if...else if" chain with identical branches (lines 17 and 23)/
		print("something")
	} else if true {
		print("something else")
	} else if true {
		print("other thing")
	} else if false {
		print("something")
	} else {
		println()
	}

	if true {
		print("something")
	} else if true {
		print("something else")
		if true { // MATCH /"if...else if" chain with identical branches (lines 33 and 35)/
			print("something")
		} else if false {
			print("something")
		} else {
			println()
		}
	}

	// Should not warn because even if branches are identical, the err variable is not
	if err := something(); err != nil {
		println(err)
	} else if err := somethingElse(); err != nil {
		println(err)
	}

	// Identical pair of branches
	if a { // MATCH /"if...else if" chain with identical branches (lines 50 and 52)/
		foo()
	} else if c {
		foo()
	} else {
		bar()
	}

	// Another identical pair of branches
	if b { // MATCH /"if...else if" chain with identical branches (lines 59 and 61)/
		bar()
	} else if d {
		bar()
	} else {
		baz()
	}

	if createFile() { // MATCH /"if...else if" chain with identical branches (lines 67 and 71)/
		doSomething()
	} else if !delete() {
		return new("cannot delete file")
	} else if createFile() {
		doSomething()
	} else {
		return new("file error")
	}

	// Test confidence is reset
	if a { // MATCH /"if...else if" chain with identical branches (lines 78 and 80)/
		foo()
	} else if b {
		foo()
	} else {
		bar()
	}
}
