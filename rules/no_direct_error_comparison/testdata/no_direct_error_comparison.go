package pkg

import (
	"errors"
	"io"
)

var (
	ErrPermission = errors.New("permission denied")
	ErrTimeout    = errors.New("timeout")
)

// Direct equality comparison of errors - should trigger
func handleErrorEq(err error) {
	if err == ErrPermission { // MATCH /avoid direct error comparison, use errors.Is(err, ErrPermission) instead/
		// handle permission error
	}
}

// Direct inequality comparison of errors - should trigger
func handleErrorNeq(err error) {
	if err != ErrTimeout { // MATCH /avoid direct error comparison, use !errors.Is(err, ErrTimeout) instead/
		// handle non-timeout error
	}
}

// Comparing errors returned from function calls - should trigger
func process() {
	err := doSomething()
	if err == getExpectedError() { // MATCH /avoid direct error comparison, use errors.Is(err, getExpectedError()) instead/
		return
	}
}

// Comparing to nil is allowed - should NOT trigger
func handleErrorNil(err error) {
	if err != nil {
		// handle error
	}
	if err == nil {
		// no error
	}
}

// Comparing to io.EOF is allowed - should NOT trigger
func readAll(buf []byte) {
	var r io.Reader
	_, err := r.Read(buf)
	if err == io.EOF {
		return
	}
	if err != io.EOF {
		return
	}
}

// Using errors.Is() - should NOT trigger
func handleErrorCorrectly(err error) {
	if errors.Is(err, ErrPermission) {
		// handle permission error
	}
	if !errors.Is(err, ErrTimeout) {
		// handle non-timeout error
	}
}

func doSomething() error         { return nil }
func getExpectedError() error    { return nil }
