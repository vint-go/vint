package fixtures

// Invalid: appending entire slice inside range loop over that same slice
func appendAllInRange() {
	var ns []int
	var rs []int
	for _, n := range ns {
		_ = n
		rs = append(rs, ns...) // MATCH /appending entire slice 'ns' inside range loop over 'ns' causes quadratic behavior/
	}
}

// Invalid: same pattern with different variable names
func appendAllInRange2() {
	var items []string
	var result []string
	for _, item := range items {
		_ = item
		result = append(result, items...) // MATCH /appending entire slice 'items' inside range loop over 'items' causes quadratic behavior/
	}
}

// Invalid: using index-only range
func appendAllIndexRange() {
	var xs []int
	var ys []int
	for range xs {
		ys = append(ys, xs...) // MATCH /appending entire slice 'xs' inside range loop over 'xs' causes quadratic behavior/
	}
}

// Valid: append just the current element
func appendSingleElement() {
	var ns []int
	var rs []int
	for _, n := range ns {
		rs = append(rs, n)
	}
}

// Valid: appending a different slice (not the one being ranged over)
func appendDifferentSlice() {
	var ns []int
	var rs []int
	other := []int{1, 2, 3}
	for _, n := range ns {
		_ = n
		rs = append(rs, other...)
	}
}

// Valid: append without ellipsis
func appendWithoutEllipsis() {
	var ns []int
	var rs []int
	for _, n := range ns {
		_ = n
		rs = append(rs, 1, 2, 3)
	}
}
