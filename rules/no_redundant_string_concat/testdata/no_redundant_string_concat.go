package fixtures

func redundantStringConcat() {
	var x string

	// Invalid: empty string concatenation on the left
	s := "" + x // MATCH /redundant concatenation with empty string/

	// Invalid: empty string concatenation on the right
	s = x + "" // MATCH /redundant concatenation with empty string/

	// Valid: normal concatenation with non-empty strings
	s = x + "hello"

	// Valid: simple assignment
	s = x

	_ = s
}
