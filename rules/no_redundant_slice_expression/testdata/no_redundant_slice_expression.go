package fixtures

// Invalid: s[:] is redundant
func redundantFullSlice() {
	s := []int{1, 2, 3}
	x := s[:] // MATCH /redundant slice expression s[:] can be simplified to s/
	_ = x
}

// Invalid: s[0:len(s)] is redundant
func redundantZeroToLen() {
	s := []int{1, 2, 3}
	x := s[0:len(s)] // MATCH /redundant slice expression s[0:len(s)] can be simplified to s/
	_ = x
}

// Invalid: string full slice
func redundantStringSlice() {
	s := "hello"
	x := s[:] // MATCH /redundant slice expression s[:] can be simplified to s/
	_ = x
}

// Valid: array[:] is a necessary conversion from [N]T to []T
func arrayToSliceConversion() {
	var a [5]int
	x := a[:]
	_ = x
}

// Valid: array[0:len(a)] is also a conversion from [N]T to []T
func arrayToSliceConversionWithLen() {
	var a [5]int
	x := a[0:len(a)]
	_ = x
}

// Valid: partial slice from start
func partialSliceFromStart() {
	s := []int{1, 2, 3}
	x := s[:2]
	_ = x
}

// Valid: partial slice from offset
func partialSliceFromOffset() {
	s := []int{1, 2, 3}
	x := s[1:]
	_ = x
}

// Valid: partial slice with both bounds
func partialSliceWithBounds() {
	s := []int{1, 2, 3}
	x := s[1:2]
	_ = x
}

// Valid: just using the variable
func directUsage() {
	s := []int{1, 2, 3}
	x := s
	_ = x
}

// Valid: s[0:len(other)] is not redundant
func zeroToLenOfDifferent() {
	s := []int{1, 2, 3}
	other := []int{1, 2}
	x := s[0:len(other)]
	_ = x
}

// Valid: s[1:len(s)] is not redundant
func oneToLen() {
	s := []int{1, 2, 3}
	x := s[1:len(s)]
	_ = x
}

// Valid: three-index slice s[:len(s):len(s)]
func threeIndexSlice() {
	s := []int{1, 2, 3}
	x := s[:len(s):len(s)]
	_ = x
}
