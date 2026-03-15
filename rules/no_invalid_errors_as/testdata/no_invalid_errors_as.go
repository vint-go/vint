package pkg

import (
	"errors"
	"fmt"
)

type MyError struct {
	Code int
}

func (e *MyError) Error() string {
	return "my error"
}

type ValueError struct {
	Msg string
}

func (e ValueError) Error() string {
	return e.Msg
}

type NotAnError struct {
	Data string
}

func example(err error) {
	var myErr MyError
	// Bad: should pass a pointer to *MyError, not a pointer to MyError
	if errors.As(err, &myErr) { // MATCH /second argument to errors.As must be a pointer to an interface or a type implementing error, got *pkg.MyError/
		fmt.Println(myErr.Code)
	}

	var notErr NotAnError
	// Bad: NotAnError does not implement error
	if errors.As(err, &notErr) { // MATCH /second argument to errors.As must be a pointer to an interface or a type implementing error, got *pkg.NotAnError/
		fmt.Println(notErr.Data)
	}

	var myErr2 *MyError
	// Good: passing a pointer to *MyError
	if errors.As(err, &myErr2) {
		fmt.Println(myErr2.Code)
	}

	var valErr *ValueError
	// Good: ValueError implements error with value receiver
	if errors.As(err, &valErr) {
		fmt.Println(valErr.Msg)
	}

	var target error
	// Good: pointer to interface
	if errors.As(err, &target) {
		fmt.Println(target)
	}

	var valErr2 ValueError
	// Good: ValueError implements error with value receiver, so *ValueError is valid
	if errors.As(err, &valErr2) {
		fmt.Println(valErr2.Msg)
	}
}
