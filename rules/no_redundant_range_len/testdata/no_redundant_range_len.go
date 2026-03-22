package fixtures

import "fmt"

func rangelenExamples() {
	s := []int{1, 2, 3}

	// Bad: unnecessary len() in range
	for i := range len(s) { // MATCH /unnecessary len() in range; use `for i := range s` instead/
		fmt.Println(i, s[i])
	}

	// Bad: unused loop variable with len()
	for _ = range len(s) { // MATCH /unnecessary len() in range; use `for range s` instead/
		fmt.Println("item")
	}

	// OK: ranging directly over the slice
	for i := range s {
		fmt.Println(i, s[i])
	}

	// OK: ranging directly without variable
	for range s {
		fmt.Println("item")
	}

	// Bad: array case
	a := [3]int{1, 2, 3}
	for i := range len(a) { // MATCH /unnecessary len() in range; use `for i := range a` instead/
		fmt.Println(i)
	}

	// OK: len() of something that is not a slice or array (e.g., map, string)
	m := map[string]int{"a": 1}
	for i := range len(m) {
		fmt.Println(i)
	}

	str := "hello"
	for i := range len(str) {
		fmt.Println(i)
	}
}
