package fixtures

import (
	"errors"
	"fmt"
)

func errorf(s string) error {
	err := errors.New(fmt.Sprintf("error: %s", s)) // MATCH /should replace errors.New(fmt.Sprintf(...)) with fmt.Errorf(...)/
	return err
}

func noErrorf(s string) error {
	err := errors.New("error: " + s)
	return err
}
