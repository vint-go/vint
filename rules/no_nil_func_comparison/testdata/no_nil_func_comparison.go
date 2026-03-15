package fixtures

import "fmt"

func myFunc() {}

func myFuncWithReturn() error {
	return nil
}

func example() {
	// Bad: comparing a named function to nil is always false
	if myFunc == nil { // MATCH /comparison of function myFunc with nil is always the same/
		fmt.Println("unreachable")
	}

	// Bad: comparing a named function to nil using !=
	if myFunc != nil { // MATCH /comparison of function myFunc with nil is always the same/
		fmt.Println("always true")
	}

	// Bad: nil on the left side
	if nil == myFunc { // MATCH /comparison of function myFunc with nil is always the same/
		fmt.Println("unreachable")
	}

	// Bad: comparing a named function with return value to nil (the function itself, not its result)
	if myFuncWithReturn == nil { // MATCH /comparison of function myFuncWithReturn with nil is always the same/
		fmt.Println("unreachable")
	}

	// Bad: comparing fmt.Println (a named function from another package) to nil
	if fmt.Println == nil { // MATCH /comparison of function fmt.Println with nil is always the same/
		fmt.Println("unreachable")
	}

	// Good: comparing a function variable to nil is meaningful
	var fn func()
	if fn == nil {
		fmt.Println("fn is not set")
	}

	// Good: comparing the result of a function call to nil
	if myFuncWithReturn() == nil {
		fmt.Println("no error")
	}

	// Good: comparing non-function values to nil
	var s *string
	if s == nil {
		fmt.Println("nil pointer")
	}
}
