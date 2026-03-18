package fixtures

// Invalid: return value of append is discarded
func discardedAppend() []int {
	s := []int{1, 2, 3}
	append(s, 4) // MATCH /result of append is discarded and has no effect/
	return s
}

// Invalid: nested append calls where the outermost result is discarded
func discardedNestedAppend() []int {
	s := []int{1, 2, 3}
	append(append(s, 4), 5) // MATCH /result of append is discarded and has no effect/
	return s
}

// Valid: result is assigned back
func validAppend() []int {
	s := []int{1, 2, 3}
	s = append(s, 4)
	return s
}

// Valid: result is used in short variable declaration
func validAppendShortDecl() []int {
	s := []int{1, 2, 3}
	s2 := append(s, 4)
	return s2
}

// Valid: result is returned directly
func validAppendReturn() []int {
	s := []int{1, 2, 3}
	return append(s, 4)
}

// Valid: result is passed to another function
func validAppendAsArg() {
	s := []int{1, 2, 3}
	_ = len(append(s, 4))
}
