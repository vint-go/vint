package fixtures

// Invalid: for-range loop appending elements one by one
func mergeSlices(a, b []int) []int {
	for _, v := range b { // MATCH /use 'a = append(a, b...)' instead of a loop to append slice elements/
		a = append(a, v)
	}
	return a
}

// Invalid: same pattern with strings
func mergeStrings(dst, src []string) []string {
	for _, s := range src { // MATCH /use 'dst = append(dst, src...)' instead of a loop to append slice elements/
		dst = append(dst, s)
	}
	return dst
}

// Valid: already using the spread operator
func mergeCorrectly(a, b []int) []int {
	a = append(a, b...)
	return a
}

// Valid: body has more than one statement
func mergeAndProcess(a, b []int) []int {
	for _, v := range b {
		a = append(a, v)
		_ = v
	}
	return a
}

// Valid: appending a transformed element
func mergeTransformed(a, b []int) []int {
	for _, v := range b {
		a = append(a, v*2)
	}
	return a
}

// Valid: different variable in append than range value
func differentVar(a, b []int) []int {
	x := 0
	for _, v := range b {
		a = append(a, x)
		_ = v
	}
	_ = x
	return a
}

// Valid: append target differs from first append arg
func mismatchedTarget(a, b, c []int) []int {
	for _, v := range b {
		a = append(c, v)
	}
	return a
}

// Valid: range value is blank identifier
func blankRange(a, b []int) []int {
	for _, _ = range b {
		a = append(a, 1)
	}
	return a
}

// Valid: no range value (index-only)
func indexOnly(a, b []int) []int {
	for i := range b {
		a = append(a, b[i])
	}
	return a
}
