package fixtures

import (
	"fmt"
	"sort"
)

func less(i, j int) bool {
	return i < j
}

func doWork() {
	fmt.Println("work")
}

func process(i, j int) {
	fmt.Println(i, j)
}

// Invalid: unnecessary lambda wrapping sort comparison
func badSortSlice() {
	xs := []int{3, 1, 2}
	sort.Slice(xs, func(i, j int) bool { // MATCH /unnecessary lambda, use the function directly/
		return less(i, j)
	})
}

// Invalid: unnecessary lambda wrapping a no-arg function in go statement
func badGoStatement() {
	go func() { doWork() }() // MATCH /unnecessary lambda, use the function directly/
}

// Invalid: unnecessary lambda wrapping a function call
func badSimpleWrapper() {
	f := func() { doWork() } // MATCH /unnecessary lambda, use the function directly/
	f()
}

// Invalid: unnecessary lambda wrapping a two-param function
func badTwoParams() {
	f := func(i, j int) { // MATCH /unnecessary lambda, use the function directly/
		process(i, j)
	}
	f(1, 2)
}

// Valid: direct function reference
func goodSortSlice() {
	xs := []int{3, 1, 2}
	sort.Slice(xs, less)
}

// Valid: direct go statement
func goodGoStatement() {
	go doWork()
}

// Valid: lambda body has more than one statement
func goodMultipleStatements() {
	f := func() {
		fmt.Println("before")
		doWork()
	}
	f()
}

// Valid: lambda body is empty
func goodEmptyBody() {
	f := func() {}
	f()
}

// Valid: arguments are not passed in order
func goodDifferentOrder() {
	f := func(i, j int) {
		process(j, i)
	}
	f(1, 2)
}

// Valid: lambda does extra work (different arguments)
func goodExtraArgs() {
	f := func(i int) {
		fmt.Println(i, "extra")
	}
	f(1)
}

// Valid: lambda body is not a call expression (assignment)
func goodNotACallExpr() {
	x := 0
	f := func() {
		x = 1
	}
	f()
	_ = x
}

// Valid: lambda has return type but body is expression statement
func goodWithReturnType() {
	f := func() error {
		return fmt.Errorf("error")
	}
	_ = f()
}

// Valid: different number of args than params
func goodDifferentArgCount() {
	f := func(i, j int) {
		fmt.Println(i)
	}
	f(1, 2)
}

// Valid: inner call uses variadic expansion
func goodVariadicExpansion() {
	f := func(args ...int) {
		fmt.Println(args...)
	}
	f(1, 2, 3)
}
