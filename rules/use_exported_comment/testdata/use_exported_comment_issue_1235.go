package golint

import (
	"errors"
)

// SendJson sends a JSON object to the server.
// MATCH /comment on exported function SendJSON should be of the form "SendJSON ..." by using its correct casing, not "SendJson ..."/
func SendJSON(data interface{}) error {
	return nil
}

// ErrInvalidJson is returned when the JSON is invalid.
// MATCH /comment on exported var ErrInvalidJSON should be of the form "ErrInvalidJSON ..." by using its correct casing, not "ErrInvalidJson ..."/
var ErrInvalidJSON = errors.New("invalid JSON")

// StatusHTTP represents an HTTP status code.
type StatusHTTP int

// Foobar blah blah
// MATCH /comment on exported method StatusHTTP.FooBar should be of the form "FooBar ..." by using its correct casing, not "Foobar ..."/
func (s StatusHTTP) FooBar() int {
	return int(s)
}

// qux was previously unexported, but now it is exported.
// MATCH /comment on exported method StatusHTTP.Qux should be of the form "Qux ..." to match its exported status, not "qux ..."/
func (s StatusHTTP) Qux() int {
	return int(s)
}

// SendJson sends a JSON object to the server.
// MATCH /comment on exported function SendJSON should be of the form "SendJSON ..." by using its correct casing, not "SendJson ..."/
func SendJSON(data interface{}) error {
	return nil
}

// errNotFound is returned when the requested resource is not found.
// MATCH /comment on exported var ErrNotFound should be of the form "ErrNotFound ..." to match its exported status, not "errNotFound ..."/
var ErrNotFound = errors.New("not found")

// VeryLongCommentThatCouldBeCJKThatCannotBeSplitOnSpaces is about the function F.
// MATCH /comment on exported function F should be of the form "F ..."/
func F() string {
	return "F"
}
