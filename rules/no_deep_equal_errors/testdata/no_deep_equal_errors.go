package fixtures

import (
	"errors"
	"io"
	"reflect"
)

func deepEqualWithErrors() {
	err1 := errors.New("something failed")
	err2 := errors.New("something failed")
	// Bad: using reflect.DeepEqual with errors
	if reflect.DeepEqual(err1, err2) { // MATCH /avoid using reflect.DeepEqual with error values, use errors.Is instead/
		// ...
	}
}

func deepEqualWithErrorInterface() {
	var err1 error = errors.New("fail")
	var err2 error = errors.New("fail")
	if reflect.DeepEqual(err1, err2) { // MATCH /avoid using reflect.DeepEqual with error values, use errors.Is instead/
		// ...
	}
}

func deepEqualOneArgIsError() {
	err := errors.New("fail")
	other := "some string"
	if reflect.DeepEqual(err, other) { // MATCH /avoid using reflect.DeepEqual with error values, use errors.Is instead/
		// ...
	}
}

func deepEqualSecondArgIsError() {
	other := "some string"
	err := errors.New("fail")
	if reflect.DeepEqual(other, err) { // MATCH /avoid using reflect.DeepEqual with error values, use errors.Is instead/
		// ...
	}
}

func deepEqualWithSentinelError() {
	var err error = io.EOF
	if reflect.DeepEqual(err, io.EOF) { // MATCH /avoid using reflect.DeepEqual with error values, use errors.Is instead/
		// ...
	}
}

// Valid: using errors.Is for error comparison
func goodErrorComparison(err error) {
	if errors.Is(err, io.EOF) {
		// ...
	}
}

// Valid: reflect.DeepEqual with non-error types
func deepEqualWithNonErrors() {
	a := []int{1, 2, 3}
	b := []int{1, 2, 3}
	if reflect.DeepEqual(a, b) {
		// ...
	}
}

// Valid: reflect.DeepEqual with strings
func deepEqualWithStrings() {
	a := "hello"
	b := "hello"
	if reflect.DeepEqual(a, b) {
		// ...
	}
}
