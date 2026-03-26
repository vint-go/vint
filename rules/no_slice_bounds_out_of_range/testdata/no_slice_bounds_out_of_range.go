package fixtures

import "sort"

// --- Should flag (no bounds check) ---

func unboundedIndex(s []int, idx int) int {
	return s[idx] // MATCH /possible slice bounds out of range/
}

func unboundedFirst(s []string) string {
	return s[0] // MATCH /possible slice bounds out of range/
}

func unboundedSubSlice(s []int) []int {
	return s[2:5] // MATCH /possible slice bounds out of range/
}

// --- Safe: len() bounds check ---

func boundedByLenCheck(s []int, idx int) int {
	if idx >= 0 && idx < len(s) {
		return s[idx]
	}
	return 0
}

func boundedNonEmpty(s []string) string {
	if len(s) > 0 {
		return s[0]
	}
	return ""
}

// --- Safe: for range loops (index matches ranged slice) ---

func rangeIndex(s []int) int {
	sum := 0
	for i := range s {
		sum += s[i]
	}
	return sum
}

func rangeIndexAddr(pes []int) *int {
	for i := range pes {
		return &pes[i]
	}
	return nil
}

// --- Unsafe: range variable used with different slice ---

func rangeMismatchSlice(a []int, b []int) {
	for i := range a {
		_ = b[i] // MATCH /possible slice bounds out of range/
	}
}

// --- Safe: sort.Slice callbacks ---

func sortSliceCallback(endpoints []string) {
	sort.Slice(endpoints, func(i, j int) bool {
		return endpoints[i] < endpoints[j]
	})
}

func sortSliceStableCallback(endpoints []string) {
	sort.SliceStable(endpoints, func(i, j int) bool {
		return endpoints[i] < endpoints[j]
	})
}

// --- Unsafe: sort.Slice parameter used with different slice ---

func sortSliceMismatch(a []string, b []string) {
	sort.Slice(a, func(i, j int) bool {
		return b[i] < a[j] // MATCH /possible slice bounds out of range/
	})
}

// --- Safe: full slice expression ---

func fullSlice(s []int) []int {
	return s[:]
}
