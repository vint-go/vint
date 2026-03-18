package fixtures

import (
	"errors"
	"fmt"
)

func sprintfWithError() string {
	err := errors.New("something went wrong")
	return fmt.Sprintf("%s", err) // MATCH /use err.Error() instead of fmt.Sprintf("%s", err)/
}

func sprintfWithErrorParam(err error) string {
	return fmt.Sprintf("%s", err) // MATCH /use err.Error() instead of fmt.Sprintf("%s", err)/
}

// Valid: using err.Error() directly
func usingErrorMethod(err error) string {
	return err.Error()
}

// Valid: format string is not just "%s"
func sprintfWithExtraFormatting(err error) string {
	return fmt.Sprintf("error: %s", err)
}

// Valid: multiple arguments
func sprintfWithMultipleArgs(err error) string {
	return fmt.Sprintf("%s: %s", "prefix", err)
}

// Valid: argument is a string, not an error
func sprintfWithString(s string) string {
	return fmt.Sprintf("%s", s)
}

// Valid: different format verb
func sprintfWithDifferentVerb(err error) string {
	return fmt.Sprintf("%v", err)
}

// Valid: no arguments beyond format string
func sprintfNoArgs() string {
	return fmt.Sprintf("hello")
}
