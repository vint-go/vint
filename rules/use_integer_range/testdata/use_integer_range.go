package fixtures

import "fmt"

func intRangeExamples() {
	n := 10
	s := []int{1, 2, 3}

	// Bad: C-style for loop that can use integer range
	for i := 0; i < 10; i++ { // MATCH /for loop can be simplified to `for i := range 10`/
		fmt.Println(i)
	}

	// Bad: unused loop variable
	for i := 0; i < 10; i++ { // MATCH /for loop can be simplified to `for range 10`/
		fmt.Println("Hello!")
	}

	// Bad: using += 1 increment
	for i := 0; i < n; i += 1 { // MATCH /for loop can be simplified to `for i := range n`/
		fmt.Println(i)
	}

	// Bad: using i = i + 1 increment
	for i := 0; i < n; i = i + 1 { // MATCH /for loop can be simplified to `for i := range n`/
		fmt.Println(i)
	}

	// Bad: using i = 1 + i increment
	for i := 0; i < n; i = 1 + i { // MATCH /for loop can be simplified to `for i := range n`/
		fmt.Println(i)
	}

	// Bad: n > i condition form
	for i := 0; n > i; i++ { // MATCH /for loop can be simplified to `for i := range n`/
		fmt.Println(i)
	}

	// Bad: using len(s) as upper bound
	for i := 0; i < len(s); i++ { // MATCH /for loop can be simplified to `for i := range len(s)`/
		fmt.Println(i)
	}

	// OK: loop variable modified in body
	for i := 0; i < 10; i++ {
		i += 2
		fmt.Println(i)
	}

	// OK: doesn't start at 0
	for i := 1; i < 10; i++ {
		fmt.Println(i)
	}

	// OK: increment by 2
	for i := 0; i < 10; i += 2 {
		fmt.Println(i)
	}

	// OK: decrement
	for i := 0; i < 10; i-- {
		fmt.Println(i)
	}

	// OK: already using range
	for i := range 10 {
		fmt.Println(i)
	}
}
