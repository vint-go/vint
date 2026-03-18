package fixtures

import "fmt"

// Invalid: empty if body with no statements.
func emptyIfBody(x int) {
	if x > 0 { // MATCH /empty branch/
	}
}

// Invalid: empty else body with no statements.
func emptyElseBody(x int) {
	if x > 0 {
		fmt.Println("positive")
	} else { // MATCH /empty branch/
	}
}

// Invalid: empty if body and empty else body.
func emptyIfAndElse(x int) {
	if x > 0 { // MATCH /empty branch/
	} else { // MATCH /empty branch/
	}
}

// Valid: if body has a statement.
func nonEmptyIfBody(x int) {
	if x > 0 {
		fmt.Println("positive")
	}
}

// Valid: else body has a statement.
func nonEmptyElseBody(x int) {
	if x > 0 {
		fmt.Println("positive")
	} else {
		fmt.Println("non-positive")
	}
}

// Valid: else-if chain, not a bare else block.
func elseIfChain(x int) {
	if x > 0 {
		fmt.Println("positive")
	} else if x < 0 {
		fmt.Println("negative")
	}
}
