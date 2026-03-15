package fixtures

import "sort"

func sortSliceInvalidPointer() {
	s := []int{3, 1, 2}
	// Bad: passing a pointer to a slice instead of the slice itself
	sort.Slice(&s, func(i, j int) bool { // MATCH /sort.Slice's first argument must be a slice; found *[]int/
		return s[i] < s[j]
	})
}

func sortSliceValid() {
	s := []int{3, 1, 2}
	// Good: passing the slice directly
	sort.Slice(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

func sortSliceStableInvalidPointer() {
	s := []string{"b", "a", "c"}
	sort.SliceStable(&s, func(i, j int) bool { // MATCH /sort.SliceStable's first argument must be a slice; found *[]string/
		return s[i] < s[j]
	})
}

func sortSliceStableValid() {
	s := []string{"b", "a", "c"}
	sort.SliceStable(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}

func sortSliceIsSortedInvalidPointer() {
	s := []int{1, 2, 3}
	sort.SliceIsSorted(&s, func(i, j int) bool { // MATCH /sort.SliceIsSorted's first argument must be a slice; found *[]int/
		return s[i] < s[j]
	})
}

func sortSliceIsSortedValid() {
	s := []int{1, 2, 3}
	sort.SliceIsSorted(s, func(i, j int) bool {
		return s[i] < s[j]
	})
}
