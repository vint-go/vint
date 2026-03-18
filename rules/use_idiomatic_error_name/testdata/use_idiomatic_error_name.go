package fixtures

import (
	"errors"
	"fmt"
)

// Invalid: sentinel errors without proper naming

var NotFoundError = errors.New("not found")       // MATCH /error var NotFoundError should have name of the form ErrFoo/
var ConnectionError = fmt.Errorf("connection lost") // MATCH /error var ConnectionError should have name of the form ErrFoo/
var timeoutError = errors.New("timeout")           // MATCH /error var timeoutError should have name of the form errFoo/
var badRequest = fmt.Errorf("bad request")         // MATCH /error var badRequest should have name of the form errFoo/

// Valid: properly named error variables

var ErrNotFound = errors.New("not found")
var ErrConnection = fmt.Errorf("connection lost")
var errTimeout = errors.New("timeout")
var errBadRequest = fmt.Errorf("bad request")
var err = errors.New("generic error")
var Err = errors.New("exported generic error")

// Valid: blank identifier
var _ = errors.New("unused")

// Valid: non-error variables (not errors.New or fmt.Errorf)
var FooBar = "not an error"
