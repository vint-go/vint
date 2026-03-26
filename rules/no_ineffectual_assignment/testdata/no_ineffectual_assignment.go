package fixtures

import "fmt"

// Invalid: Variable assigned but immediately overwritten
func exampleOverwritten() int {
	x := 1 // MATCH /ineffectual assignment to x/
	x = 2
	return x
}

// Invalid: Error variable assigned but overwritten without checking
func exampleErrorOverwritten() error {
	err := doSomething()    // MATCH /ineffectual assignment to err/
	err = doSomethingElse()
	return err
}

// Invalid: Variable assigned but shadowed by reassignment
func exampleReassigned() {
	x := 1
	fmt.Println(x)
	x = 2 // MATCH /ineffectual assignment to x/
	x = 3
	fmt.Println(x)
}

// Invalid: Variable assigned in all branches before being read
func exampleAllBranches(flag bool) int {
	x := 0 // MATCH /ineffectual assignment to x/
	if flag {
		x = 1
	} else {
		x = 2
	}
	return x
}

// Valid: Variable assigned and subsequently used
func exampleValid() int {
	x := 1
	return x
}

// Valid: Variable assigned and read before reassignment
func exampleValidReadBeforeReassign() {
	x := computeValue()
	fmt.Println(x)
	x = otherValue()
	fmt.Println(x)
}

// Valid: Error variable checked before reassignment
func exampleValidErrorChecked() error {
	err := doSomething()
	if err != nil {
		return err
	}
	err = doSomethingElse()
	return err
}

// Valid: Variable assigned and used in a loop
func exampleValidLoop(items []int) int {
	sum := 0
	for _, item := range items {
		sum += item
	}
	return sum
}

// Valid: Blank identifier used intentionally
func exampleValidBlankIdent() {
	_, err := doSomethingMulti()
	if err != nil {
		fmt.Println(err)
	}
}

// Valid: Variable used in condition
func exampleValidCondition() {
	x := 5
	if x > 3 {
		fmt.Println(x)
	}
}

// Valid: Variable used after if/else (not touched in branches) - Bug 1 regression test
func exampleValidUsedAfterIfElse(flag bool) int {
	x := 10
	if flag {
		fmt.Println("branch a")
	} else {
		fmt.Println("branch b")
	}
	return x
}

// Valid: Initial value used in loop - Bug 1 regression test
func exampleValidInitialValueInLoop() int {
	failedAmount := 0
	for i := 0; i < 5; i++ {
		delay := 100 + failedAmount*failedAmount*100
		fmt.Println(delay)
		failedAmount++
	}
	return failedAmount
}

// Valid: Read-then-reassign via append - Bug 2 regression test
func exampleValidAppendReadReassign() []int {
	result := []int{1, 2, 3}
	if condition {
		result = append(result, 4)
	}
	return result
}

// Valid: Read-then-reassign in both branches via append - Bug 2 regression test
func exampleValidAppendBothBranches(flag bool) []int {
	items := []int{1}
	if flag {
		items = append(items, 2)
	} else {
		items = append(items, 3)
	}
	return items
}

// Helper functions for the test fixture to parse correctly
func doSomething() error             { return nil }
func doSomethingElse() error         { return nil }
func computeValue() int              { return 0 }
func otherValue() int                { return 0 }
func doSomethingMulti() (int, error) { return 0, nil }

var condition bool
