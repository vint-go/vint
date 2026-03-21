package example

import "errors"

func validate(name string) error {
	if name == "" {
		return errors.New("name is required") //nolint:err113
	}
	return nil
}
