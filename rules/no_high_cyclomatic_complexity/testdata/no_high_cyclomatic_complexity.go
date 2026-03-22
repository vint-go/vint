// Test of noHighCyclomaticComplexity rule.

// Package pkg ...
package pkg

import (
	"errors"
	"fmt"
)

// Low cyclomatic complexity - clean, focused function (complexity = 1)
func add(a, b int) int {
	return a + b
}

// Moderate complexity within acceptable range (complexity = 2)
func greet(name string, formal bool) string {
	if formal {
		return "Hello, " + name + "."
	}
	return "Hi, " + name + "!"
}

// complexity = 1 + 5 ifs = 6, exceeds threshold of 5
func highComplexity(a, b, c, d, e int) int { // MATCH /function highComplexity has cyclomatic complexity 6 (> max enabled 5)/
	if a > 0 {
		a++
	}
	if b > 0 {
		b++
	}
	if c > 0 {
		c++
	}
	if d > 0 {
		d++
	}
	if e > 0 {
		e++
	}
	return a + b + c + d + e
}

// complexity = 1 (func) + 8 cases + 3 ifs + 1 for = 13, exceeds threshold of 5 (default case excluded)
func evaluate(op string, a, b int) (int, error) { // MATCH /function evaluate has cyclomatic complexity 13 (> max enabled 5)/
	switch op {
	case "add":
		return a + b, nil
	case "sub":
		return a - b, nil
	case "mul":
		return a * b, nil
	case "div":
		if b == 0 {
			return 0, errors.New("division by zero")
		}
		return a / b, nil
	case "mod":
		if b == 0 {
			return 0, errors.New("modulo by zero")
		}
		return a % b, nil
	case "pow":
		result := 1
		for i := 0; i < b; i++ {
			result *= a
		}
		return result, nil
	case "max":
		if a > b {
			return a, nil
		}
		return b, nil
	case "avg":
		return (a + b) / 2, nil
	default:
		return 0, fmt.Errorf("unknown operation: %s", op)
	}
}

// Function excluded from complexity check using ignore directive
//
//gocyclo:ignore
func complexButIgnored(a, b, c, d, e, f int) int {
	if a > 0 {
		a++
	}
	if b > 0 {
		b++
	}
	if c > 0 {
		c++
	}
	if d > 0 {
		d++
	}
	if e > 0 {
		e++
	}
	if f > 0 {
		f++
	}
	return a + b + c + d + e + f
}

// complexity = 1 + 2 ifs + 1 && + 1 || = 5, does NOT exceed threshold of 5
func moderateComplexity(x, y int) bool {
	if x > 0 && y > 0 {
		return true
	}
	if x < 0 || y < 0 {
		return false
	}
	return x == y
}
