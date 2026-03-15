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

// Helper functions for the test fixture to parse correctly
func doSomething() error            { return nil }
func doSomethingElse() error        { return nil }
func computeValue() int             { return 0 }
func otherValue() int               { return 0 }
func doSomethingMulti() (int, error) { return 0, nil }

var condition bool
