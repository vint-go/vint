package fixtures

import (
	"errors"
	"fmt"
)

type MyError struct {
	Code int
}

func (e *MyError) Error() string {
	return fmt.Sprintf("code: %d", e.Code)
}

type OtherError struct {
	Message string
}

func (e *OtherError) Error() string {
	return e.Message
}

func badTypeAssert(err error) {
	myErr, ok := err.(*MyError) // MATCH /type assertion on error will fail on wrapped errors, use errors.As/
	if ok {
		fmt.Println(myErr.Code)
	}
}

func badTypeSwitch(err error) {
	switch e := err.(type) { // MATCH /type switch on error will fail on wrapped errors, use errors.As/
	case *MyError:
		fmt.Println(e.Code)
	case *OtherError:
		fmt.Println(e.Message)
	}
}

func goodErrorsAs(err error) {
	var myErr *MyError
	if errors.As(err, &myErr) {
		fmt.Println(myErr.Code)
	}
}

func goodMultipleErrorsAs(err error) {
	var myErr *MyError
	var otherErr *OtherError
	if errors.As(err, &myErr) {
		fmt.Println(myErr.Code)
	} else if errors.As(err, &otherErr) {
		fmt.Println(otherErr.Message)
	}
}

type notAnError struct {
	value int
}

func typeAssertOnNonError() {
	var x interface{} = &notAnError{value: 42}
	if v, ok := x.(*notAnError); ok {
		fmt.Println(v.value)
	}
}
