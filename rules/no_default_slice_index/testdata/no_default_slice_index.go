package fixtures

// Invalid: s[:len(s)] can omit the high index
func defaultSliceHigh() {
	s := []int{1, 2, 3}
	x := s[:len(s)] // MATCH /omit default slice index: s[:len(s)] can be simplified/
	_ = x
}

// Invalid: with a non-nil low index, s[1:len(s)] can be simplified to s[1:]
func defaultSliceHighWithLow() {
	s := []int{1, 2, 3}
	x := s[1:len(s)] // MATCH /omit default slice index: s[1:len(s)] can be simplified/
	_ = x
}

// Invalid: string slice
func defaultStringSliceHigh() {
	s := "hello"
	x := s[:len(s)] // MATCH /omit default slice index: s[:len(s)] can be simplified/
	_ = x
}

// Valid: different variable in len
func differentLenArg() {
	s := []int{1, 2, 3}
	other := []int{1, 2}
	x := s[:len(other)]
	_ = x
}

// Valid: using a numeric literal as high index
func numericHighIndex() {
	s := []int{1, 2, 3}
	x := s[:2]
	_ = x
}

// Valid: no high index (just s[:])
func noHighIndex() {
	s := []int{1, 2, 3}
	x := s[:]
	_ = x
}

// Valid: three-index slice expression
func threeIndexSlice() {
	s := []int{1, 2, 3}
	x := s[:len(s):len(s)]
	_ = x
}

// Valid: cap instead of len
func capHighIndex() {
	s := []int{1, 2, 3}
	x := s[:cap(s)]
	_ = x
}

// Valid: using the variable directly
func directUsage() {
	s := []int{1, 2, 3}
	x := s
	_ = x
}
