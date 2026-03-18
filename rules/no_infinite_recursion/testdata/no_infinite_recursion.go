package fixtures

import (
	"fmt"
	"os"
)

// --- Invalid cases ---

// Simple direct recursion with no base case
func factorial(n int) int {
	return n * factorial(n-1) // MATCH /infinite recursive call/
}

// Method recursion with no base case
type MyStruct struct {
	val int
}

func (s *MyStruct) recurse() {
	s.recurse() // MATCH /infinite recursive call/
}

// Function that calls itself in all branches
func allPaths(x int) int {
	if x > 0 {
		fmt.Println("positive")
	}
	return allPaths(x - 1) // MATCH /infinite recursive call/
}

// Recursion inside unconditional for loop
func loopRecursion() {
	for {
		loopRecursion() // MATCH /infinite recursive call/
		break
	}
}

// --- Valid cases ---

// Proper recursion with base case
func factorialOk(n int) int {
	if n <= 1 {
		return 1
	}
	return n * factorialOk(n-1)
}

// Method recursion with a conditional return
func (s *MyStruct) recurseOk() {
	if s.val <= 0 {
		return
	}
	s.val--
	s.recurseOk()
}

// Recursion guarded by panic in a branch
func withPanic(n int) int {
	if n < 0 {
		panic("negative")
	}
	return withPanic(n - 1)
}

// Recursion guarded by os.Exit
func withExit(n int) {
	if n == 0 {
		os.Exit(0)
	}
	withExit(n - 1)
}

// Non-recursive function
func notRecursive(n int) int {
	return n + 1
}

// Calling a different function (not self)
func a(n int) int {
	return b(n)
}

func b(n int) int {
	return n
}

// Recursion in a conditional for loop is fine (it's inside conditional)
func condForLoop(n int) {
	for n > 0 {
		condForLoop(n - 1)
		n--
	}
}

// Recursion inside if-else both branches is conditional
func inIfElse(n int) int {
	if n > 0 {
		return inIfElse(n - 1)
	}
	return 0
}

// Recursion inside switch is conditional
func inSwitch(n int) int {
	switch {
	case n > 0:
		return inSwitch(n - 1)
	default:
		return 0
	}
}
