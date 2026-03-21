package example

import "errors"

var ErrNotFound = errors.New("not found")

func lookup(key string) error {
	err := find(key)
	if err == ErrNotFound { //nolint:errorlint
		return nil
	}
	return err
}

func find(key string) error { return nil }
