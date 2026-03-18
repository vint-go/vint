package fixtures

import "sort"

func sortSliceExamples() {
	xs := []int{3, 1, 2}
	ys := []int{6, 4, 5}

	// Invalid: comparing ys instead of xs (missing slice reference)
	sort.Slice(xs, func(i, j int) bool { // MATCH /sort.Slice's comparison function does not reference the sorted slice xs/
		return ys[i] < ys[j]
	})

	// Invalid: reversed indices
	sort.Slice(xs, func(i, j int) bool { // MATCH /sort.Slice's comparison function has reversed indices: xs[j] compared with xs[i]/
		return xs[j] < xs[i]
	})

	// Valid: correct usage
	sort.Slice(xs, func(i, j int) bool {
		return xs[i] < xs[j]
	})

	// Valid: SliceStable correct usage
	sort.SliceStable(xs, func(i, j int) bool {
		return xs[i] < xs[j]
	})

	// Invalid: SliceStable with wrong slice reference
	sort.SliceStable(xs, func(i, j int) bool { // MATCH /sort.Slice's comparison function does not reference the sorted slice xs/
		return ys[i] < ys[j]
	})

	// Invalid: SliceStable with reversed indices
	sort.SliceStable(xs, func(i, j int) bool { // MATCH /sort.Slice's comparison function has reversed indices: xs[j] compared with xs[i]/
		return xs[j] < xs[i]
	})
}
