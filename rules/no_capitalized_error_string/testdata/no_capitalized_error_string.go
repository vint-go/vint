package fixtures

import (
	"errors"
	"fmt"
)

// Invalid: capitalized error string
func bad1() error {
	return errors.New("Something went wrong") // MATCH /error strings should not be capitalized/
}

// Invalid: error string ending with period
func bad2() error {
	return errors.New("something went wrong.") // MATCH /error strings should not end with punctuation or a newline/
}

// Invalid: error string ending with exclamation mark
func bad3() error {
	return errors.New("something went wrong!") // MATCH /error strings should not end with punctuation or a newline/
}

// Invalid: error string ending with colon
func bad4() error {
	return errors.New("something went wrong:") // MATCH /error strings should not end with punctuation or a newline/
}

// Invalid: capitalized error string in fmt.Errorf
func bad5() error {
	return fmt.Errorf("Something went wrong: %v", "detail") // MATCH /error strings should not be capitalized/
}

// Invalid: fmt.Errorf ending with period
func bad6() error {
	return fmt.Errorf("something went wrong.") // MATCH /error strings should not end with punctuation or a newline/
}

// Valid: lowercase error string
func good1() error {
	return errors.New("something went wrong")
}

// Valid: error string starting with proper noun / acronym (e.g. HTTP, URL, GitHub)
func good2() error {
	return errors.New("HTTP request failed")
}

// Valid: error string starting with acronym
func good3() error {
	return errors.New("URL is invalid")
}

// Valid: error string starting with Go exported identifier
func good4() error {
	return errors.New("EOF reached")
}

// Valid: error string starting with I2000-style word
func good5() error {
	return errors.New("I2000 error occurred")
}

// Valid: empty error string
func good6() error {
	return errors.New("")
}

// Valid: fmt.Errorf with lowercase
func good7() error {
	return fmt.Errorf("failed to process: %v", "x")
}

// Valid: non-string argument (dynamic)
func good8(msg string) error {
	return errors.New(msg)
}

// Valid: GitHub starts with uppercase but has more uppercase letters
func good9() error {
	return errors.New("GitHub integration failed")
}
