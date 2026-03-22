package fixtures

import (
	"errors"
	"fmt"
)

func nonWrappingV(err error) error {
	return fmt.Errorf("failed to process: %v", err) // MATCH /non-wrapping format verb for fmt.Errorf; use %w to format error err/
}

func nonWrappingS(err error) error {
	return fmt.Errorf("operation failed: %s", err) // MATCH /non-wrapping format verb for fmt.Errorf; use %w to format error err/
}

func wrappingW(err error) error {
	return fmt.Errorf("failed to process: %w", err)
}

func multipleW(err1, err2 error) error {
	return fmt.Errorf("error1: %w, error2: %w", err1, err2)
}

func nonErrorArg(s string) error {
	return fmt.Errorf("failed: %v", s)
}

func mixedArgs(s string, err error) error {
	return fmt.Errorf("context: %s, error: %v", s, err) // MATCH /non-wrapping format verb for fmt.Errorf; use %w to format error err/
}

func noArgs() error {
	return fmt.Errorf("static error")
}

func intArg(n int) error {
	return fmt.Errorf("code: %d", n)
}

func wrappedWithContext(err error) error {
	return fmt.Errorf("context: %w", err)
}

func errorsNew() error {
	return errors.New("plain error")
}
