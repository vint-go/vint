package fixtures

import (
	"fmt"
	"time"
)

// Invalid: Magic number 3 in case clause.
func categorize(x interface{}) {
	switch x {
	case "test":
	case 3: // MATCH /magic number: 3, in <case> detected/
	}
}

// Invalid: Magic numbers in expressionless switch with binary comparisons.
func greet() {
	t := time.Now()
	switch {
	case t.Hour() < 12: // MATCH /magic number: 12, in <case> detected/
		fmt.Println("Good morning!")
	case 17 > t.Hour(): // MATCH /magic number: 17, in <case> detected/
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}
}

// Invalid: Magic float 3.0 in case clause.
func categorizeFloat(x interface{}) {
	switch x {
	case 1.0:
	case 0.0:
	case 3.0: // MATCH /magic number: 3, in <case> detected/
	}
}

// Valid: Using a named constant as case value.
const categoryThreshold = 3

func categorizeGood(x interface{}) {
	switch x {
	case "test":
	case categoryThreshold:
	}
}

// Valid: Using named constants for time comparisons.
const (
	morningEnd   = 12
	afternoonEnd = 17
)

func greetGood() {
	t := time.Now()
	switch {
	case t.Hour() < morningEnd:
		fmt.Println("Good morning!")
	case afternoonEnd > t.Hour():
		fmt.Println("Good afternoon.")
	default:
		fmt.Println("Good evening.")
	}
}

// Valid: 1.0 and 0.0 are excluded by default.
func categorizeDefault(x interface{}) {
	switch x {
	case 1.0:
	case 0.0:
	}
}

// Valid: 0 and 1 are excluded by default.
func categorizeSmall(x interface{}) {
	switch x {
	case 0:
	case 1:
	}
}
