package fixtures

import "errors"

func uniqueValidation(x int) error {
	if x < 0 {
		return errors.New("must be non-negative")
	}
	return nil
}
