package fixtures

import "fmt"

type MyError struct{}

func (e *MyError) Error() string { return "error" }

type myStringer struct {
	val string
}

func (s *myStringer) String() string { return s.val }

// --- Helper functions for testing ---

// neverReturnsNil: all return paths return concrete types, never untyped nil.
func neverReturnsNil() error {
	var err *MyError
	return err
}

// alwaysConcreteReturn: returns a literal concrete error.
func alwaysConcreteReturn() error {
	return &MyError{}
}

// multiReturnNeverNil: the error position always returns a concrete type.
func multiReturnNeverNil() (int, error) {
	var err *MyError
	return 0, err
}

// canReturnNil: has a path that returns untyped nil.
func canReturnNil(flag bool) error {
	if flag {
		return nil
	}
	return &MyError{}
}

// returnsInterfaceVar: returns an interface-typed variable (could be nil).
func returnsInterfaceVar(err error) error {
	return err
}

// concreteReturnType: signature returns concrete type, not interface.
func concreteReturnType() *MyError {
	return &MyError{}
}

// neverNilStringer: never returns nil fmt.Stringer.
func neverNilStringer() fmt.Stringer {
	return &myStringer{val: "hello"}
}

// canReturnNilStringer: can return nil.
func canReturnNilStringer(flag bool) fmt.Stringer {
	if flag {
		return nil
	}
	return &myStringer{val: "hello"}
}

// mixedReturns: one position is never-nil, the other can be nil.
func mixedReturns(flag bool) (fmt.Stringer, error) {
	if flag {
		return &myStringer{val: "a"}, nil
	}
	return &myStringer{val: "b"}, &MyError{}
}

// namedResultCanBeNil: named result with bare return - conservatively may be nil.
func namedResultCanBeNil() (err error) {
	return
}

// --- Test cases: Should flag (impossible nil comparison) ---

// Bad: comparing result of function that never returns nil.
func testNeverNilEqual() {
	err := neverReturnsNil()
	if err == nil { // MATCH /nil comparison of err is always false because neverReturnsNil never returns nil/
		fmt.Println("unreachable")
	}
}

// Bad: != nil with never-nil function.
func testNeverNilNotEqual() {
	err := neverReturnsNil()
	if err != nil { // MATCH /nil comparison of err is always true because neverReturnsNil never returns nil/
		fmt.Println("always reached")
	}
}

// Bad: function that always returns concrete literal.
func testAlwaysConcreteReturn() {
	err := alwaysConcreteReturn()
	if err == nil { // MATCH /nil comparison of err is always false because alwaysConcreteReturn never returns nil/
		fmt.Println("unreachable")
	}
}

// Bad: multi-return where error position is never nil.
func testMultiReturnNeverNil() {
	_, err := multiReturnNeverNil()
	if err == nil { // MATCH /nil comparison of err is always false because multiReturnNeverNil never returns nil/
		fmt.Println("unreachable")
	}
}

// Bad: Stringer that never returns nil.
func testNeverNilStringerComparison() {
	s := neverNilStringer()
	if s == nil { // MATCH /nil comparison of s is always false because neverNilStringer never returns nil/
		fmt.Println("unreachable")
	}
}

// Bad: concrete return type assigned to interface variable.
func testConcreteToInterface() {
	var err error
	err = concreteReturnType()
	if err == nil { // MATCH /nil comparison of err is always false because concreteReturnType returns concrete type *fixtures.MyError/
		fmt.Println("unreachable")
	}
}

// Bad: concrete return type via var declaration.
func testConcreteToInterfaceVarDecl() {
	var err error = concreteReturnType()
	if err == nil { // MATCH /nil comparison of err is always false because concreteReturnType returns concrete type *fixtures.MyError/
		fmt.Println("unreachable")
	}
}

// Bad: mixed returns - Stringer position is never nil, error can be nil.
func testMixedReturnsNeverNil() {
	s, _ := mixedReturns(true)
	if s == nil { // MATCH /nil comparison of s is always false because mixedReturns never returns nil/
		fmt.Println("unreachable")
	}
}

// --- Test cases: Should NOT flag ---

// Good: function that can return nil.
func testCanReturnNilFunc() {
	err := canReturnNil(true)
	if err == nil {
		fmt.Println("valid")
	}
}

// Good: function that returns interface variable (could be nil).
func testReturnsInterfaceVar() {
	err := returnsInterfaceVar(nil)
	if err == nil {
		fmt.Println("valid")
	}
}

// Good: concrete pointer nil check (not interface comparison).
func testConcretePointerNil() {
	var err *MyError
	if err == nil {
		fmt.Println("valid - pointer nil check, not interface nil check")
	}
}

// Good: interface variable with no function call assignment.
func testInterfaceVarNoCall() {
	var err error
	if err == nil {
		fmt.Println("valid - no function call source")
	}
}

// Good: variable reassigned after initial function call.
func testReassigned() {
	err := neverReturnsNil()
	err = canReturnNil(true)
	if err == nil {
		fmt.Println("valid - err was reassigned")
	}
	_ = err
}

// Good: concrete return type assigned to concrete variable (not interface).
func testConcreteReturnAsConcreteVar() {
	err := concreteReturnType()
	if err == nil {
		fmt.Println("valid - concrete variable, pointer nil check")
	}
}

// Good: Stringer that can return nil.
func testCanReturnNilStringer() {
	s := canReturnNilStringer(true)
	if s == nil {
		fmt.Println("valid")
	}
}

// Good: named result with bare return (conservatively may be nil).
func testNamedResultBareReturn() {
	err := namedResultCanBeNil()
	if err == nil {
		fmt.Println("valid - bare return could leave named result as nil")
	}
}

// Good: mixed returns - error position CAN be nil.
func testMixedReturnsCanBeNil() {
	_, err := mixedReturns(true)
	if err == nil {
		fmt.Println("valid - error position can be nil")
	}
}

// Good: non-interface comparison in nil check context.
func testNonInterfaceReturn() {
	p := concreteReturnType()
	if p == nil {
		fmt.Println("valid - comparing concrete pointer")
	}
}
