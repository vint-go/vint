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

// --- Safe: cross-slice indexing with make([]T, len(rangedSlice)) ---

func rangeCrossSliceMakeLen(pes []int) []int {
	result := make([]int, len(pes))
	for i := range pes {
		result[i] = pes[i] * 2 // safe: result was make'd with len(pes)
	}
	return result
}

// --- Safe: range over make'd slice, index original ---

func rangeMakeLenReverse(items []string) []int {
	lengths := make([]int, len(items))
	for i := range lengths {
		lengths[i] = len(items[i]) // safe: lengths was make'd with len(items)
	}
	return lengths
}

// --- Safe: C-style for loop with make([]T, n) ---

func cStyleForWithMake(length int) []byte {
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = byte(i) // safe: result was make'd with length, loop bounded by length
	}
	return result
}

// --- Safe: C-style for loop with make([]T, len(other)) ---

func cStyleForWithMakeLenOther(src []int) []int {
	dst := make([]int, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = src[i] * 2 // safe: dst was make'd with len(src), loop bounded by len(src)
	}
	return dst
}

// --- Unsafe: C-style for loop, different bound than make ---

func cStyleForMismatch(n int, m int) []byte {
	result := make([]byte, n)
	for i := 0; i < m; i++ {
		result[i] = byte(i) // MATCH /possible slice bounds out of range/
	}
	return result
}

// --- Unsafe: cross-slice in range, no make relationship ---

func rangeCrossSliceNoMake(a []int, b []int) {
	for i := range a {
		_ = b[i] // MATCH /possible slice bounds out of range/
	}
}

// --- Safe: copier.Copy makes slices same length ---

// mock copier for testing (no external dependency needed)
var copier = struct {
	Copy func(dst, src interface{}) error
}{}

func copierCopyRangeSource(endpoints []string) []string {
	var endpointsResponse []string
	_ = copier.Copy(&endpointsResponse, endpoints)
	for i := range endpointsResponse {
		_ = endpoints[i] // safe: copier.Copy paired them
	}
	return endpointsResponse
}

func copierCopyRangeDest(capabilities []int) []int {
	var entity []int
	_ = copier.Copy(&entity, capabilities)
	for i := range capabilities {
		entity[i] = capabilities[i] * 2 // safe: copier.Copy paired them
	}
	return entity
}
