package fixtures

// Invalid: for i, v := range src { dst[i] = v }
func copyWithRangeValue() {
	src := []int{1, 2, 3}
	dst := make([]int, len(src))
	for i, v := range src { // MATCH /should replace loop with copy(dst, src)/
		dst[i] = v
	}
	_ = dst
}

// Invalid: for i := range src { dst[i] = src[i] }
func copyWithIndexOnly() {
	src := []int{1, 2, 3}
	dst := make([]int, len(src))
	for i := range src { // MATCH /should replace loop with copy(dst, src)/
		dst[i] = src[i]
	}
	_ = dst
}

// Invalid: for i, _ := range src { dst[i] = src[i] }
func copyWithBlankValue() {
	src := []int{1, 2, 3}
	dst := make([]int, len(src))
	for i, _ := range src { // MATCH /should replace loop with copy(dst, src)/
		dst[i] = src[i]
	}
	_ = dst
}

// Valid: loop body has more than one statement
func copyWithExtraWork() {
	src := []int{1, 2, 3}
	dst := make([]int, len(src))
	for i, v := range src {
		dst[i] = v
		_ = i
	}
	_ = dst
}

// Valid: assignment uses different index
func copyWithDifferentIndex() {
	src := []int{1, 2, 3}
	dst := make([]int, len(src))
	j := 0
	for _, v := range src {
		dst[j] = v
		j++
	}
	_ = dst
}

// Valid: RHS is not the range value variable
func copyWithTransformation() {
	src := []int{1, 2, 3}
	dst := make([]int, len(src))
	for i, v := range src {
		dst[i] = v + 1
	}
	_ = dst
}

// Valid: proper use of copy builtin
func properCopy() {
	src := []int{1, 2, 3}
	dst := make([]int, len(src))
	copy(dst, src)
	_ = dst
}

// Valid: for i := range src { dst[i] = other[i] } - different source
func copyFromDifferentSource() {
	src := []int{1, 2, 3}
	other := []int{4, 5, 6}
	dst := make([]int, len(src))
	for i := range src {
		dst[i] = other[i]
	}
	_ = dst
}
