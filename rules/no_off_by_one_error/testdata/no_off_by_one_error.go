package fixtures

// Invalid: accessing one past the last element.
func accessPastEnd(s []int) {
	last := s[len(s)] // MATCH /off-by-one error: index equals length of the container/
	_ = last
}

// Invalid: accessing one past the last element of a string.
func accessPastEndString(s string) {
	ch := s[len(s)] // MATCH /off-by-one error: index equals length of the container/
	_ = ch
}

// Invalid: field selector expression.
type myStruct struct {
	items []int
}

func accessPastEndField(m myStruct) {
	last := m.items[len(m.items)] // MATCH /off-by-one error: index equals length of the container/
	_ = last
}

// Valid: correct last element access.
func correctLastElement(s []int) {
	last := s[len(s)-1]
	_ = last
}

// Valid: indexing with a variable.
func indexWithVar(s []int, i int) {
	v := s[i]
	_ = v
}

// Valid: len of a different slice.
func lenOfDifferent(s []int, t []int) {
	v := s[len(t)]
	_ = v
}

// Valid: using len in a loop bound (not as direct index).
func loopBound(s []int) {
	for i := 0; i < len(s); i++ {
		_ = s[i]
	}
}
