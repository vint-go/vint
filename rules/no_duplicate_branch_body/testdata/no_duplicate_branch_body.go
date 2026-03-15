package fixtures

import "fmt"

func duplicateBranchBodies(condition bool) {
	if condition { // MATCH /both branches of the if statement have identical bodies/
		fmt.Println("result")
	} else {
		fmt.Println("result")
	}
}

func differentBranchBodies(condition bool) {
	fmt.Println("true result")
	if condition {
		fmt.Println("true result")
	} else {
		fmt.Println("false result")
	}
}

func ifWithoutElse(condition bool) {
	if condition {
		fmt.Println("only then")
	}
}

func noBranch() {
	fmt.Println("result")
}

func duplicateMultipleStatements(condition bool) {
	if condition { // MATCH /both branches of the if statement have identical bodies/
		x := 1
		fmt.Println(x)
		fmt.Println("done")
	} else {
		x := 1
		fmt.Println(x)
		fmt.Println("done")
	}
}

func differentStatementCounts(condition bool) {
	if condition {
		fmt.Println("a")
		fmt.Println("b")
	} else {
		fmt.Println("a")
	}
}

func elseIfChain(a, b bool) {
	// else-if chains: the outer if-else is not a simple if-else so it's skipped,
	// but the inner if-else has different bodies so it doesn't match either
	if a {
		fmt.Println("a")
	} else if b {
		fmt.Println("b")
	} else {
		fmt.Println("c")
	}
}
