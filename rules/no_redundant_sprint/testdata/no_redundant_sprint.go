package fixtures

import "fmt"

func redundantSprintExamples() {
	s := "hello"
	_ = fmt.Sprint(s) // MATCH /redundant fmt.Sprint call on a value that is already a string/

	str := "world"
	_ = fmt.Sprintf("%s", str) // MATCH /redundant fmt.Sprintf call with %s on a value that is already a string/
}

func validSprintExamples() {
	s := "hello"
	_ = s // direct use, no sprint

	str := "value"
	_ = fmt.Sprintf("prefix: %s", str) // not redundant, has additional formatting

	n := 42
	_ = fmt.Sprint(n) // not redundant, n is an int

	_ = fmt.Sprintf("%d", n) // not redundant, different verb

	_ = fmt.Sprint("a", "b") // not redundant, multiple arguments

	_ = fmt.Sprintf("%s %s", "a", "b") // not redundant, multiple arguments
}

type myString string

func namedStringType() {
	var ms myString = "hello"
	_ = fmt.Sprint(ms)        // MATCH /redundant fmt.Sprint call on a value that is already a string/
	_ = fmt.Sprintf("%s", ms) // MATCH /redundant fmt.Sprintf call with %s on a value that is already a string/
}
