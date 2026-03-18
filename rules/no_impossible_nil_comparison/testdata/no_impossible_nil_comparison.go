package fixtures

import "fmt"

type MyError struct{}

func (e *MyError) Error() string { return "error" }

// Bad: returning concrete *MyError as error interface - the interface will
// never be nil even when err is nil.
func getError() error {
	var err *MyError
	return err // MATCH /returning concrete type *fixtures.MyError as interface will produce non-nil interface value/
}

// Bad: returning concrete type variable in a function with interface return.
func getErrorFromParam(err *MyError) error {
	return err // MATCH /returning concrete type *fixtures.MyError as interface will produce non-nil interface value/
}

// Bad: returning concrete type in a multi-return function.
func getErrorMulti() (int, error) {
	var err *MyError
	return 0, err // MATCH /returning concrete type *fixtures.MyError as interface will produce non-nil interface value/
}

// Good: returning untyped nil is fine.
func getErrorNil() error {
	return nil
}

// Good: returning an interface-typed variable is fine.
func getErrorIface(err error) error {
	return err
}

// Good: returning concrete value after checking for nil and returning explicit nil.
func getErrorSafe() error {
	var err *MyError
	if err == nil {
		return nil
	}
	return err // MATCH /returning concrete type *fixtures.MyError as interface will produce non-nil interface value/
}

// Good: returning from a function that returns a concrete type (not an interface).
func getPointer() *MyError {
	var err *MyError
	return err
}

// Good: function returning non-interface types only.
func getInt() int {
	return 42
}

// Bad: returning a new concrete struct as an interface.
func getNewError() error {
	return &MyError{} // MATCH /returning concrete type *fixtures.MyError as interface will produce non-nil interface value/
}

type myStringer struct {
	val string
}

func (s *myStringer) String() string { return s.val }

// Bad: returning concrete type as fmt.Stringer interface.
func getStringer() fmt.Stringer {
	var s *myStringer
	return s // MATCH /returning concrete type *fixtures.myStringer as interface will produce non-nil interface value/
}

// Good: returning explicit nil as fmt.Stringer.
func getStringerNil() fmt.Stringer {
	return nil
}

// Bad: nested function literal returning concrete type as interface.
func closureExample() {
	f := func() error {
		var err *MyError
		return err // MATCH /returning concrete type *fixtures.MyError as interface will produce non-nil interface value/
	}
	_ = f
}
