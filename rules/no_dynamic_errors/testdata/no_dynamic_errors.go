package fixtures

import (
	"errors"
	"fmt"
	"os"
)

// Package-level sentinel errors are allowed
var (
	ErrNameRequired   = errors.New("name is required")
	ErrDivisionByZero = errors.New("division by zero")
)

var ErrNotFound = errors.New("not found")
var ErrUnauthorized = errors.New("unauthorized")

// Invalid: errors.New() inside a function
func validate(name string) error {
	if name == "" {
		return errors.New("name is required") // MATCH /use of errors.New() inside a function: define sentinel errors at package level/
	}
	return nil
}

// Invalid: fmt.Errorf() without %w inside a function
func openFile(path string) error {
	_, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open %s: %v", path, err) // MATCH /use of fmt.Errorf() without %w: use fmt.Errorf with %w to wrap sentinel errors/
	}
	return nil
}

// Invalid: errors.New() in a return statement
func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division by zero") // MATCH /use of errors.New() inside a function: define sentinel errors at package level/
	}
	return a / b, nil
}

// Valid: returning a package-level sentinel error
func validateOk(name string) error {
	if name == "" {
		return ErrNameRequired
	}
	return nil
}

// Valid: wrapping a sentinel error with fmt.Errorf using %w
var ErrOpenFile = errors.New("failed to open file")

func openFileOk(path string) error {
	_, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	return nil
}

// Valid: wrapping errors with %w inside functions
func process(id int) error {
	_, err := os.Open("test")
	if err != nil {
		return fmt.Errorf("processing item %d: %w", id, err)
	}
	return nil
}
