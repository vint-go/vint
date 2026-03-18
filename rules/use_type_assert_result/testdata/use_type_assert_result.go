package fixtures

import "fmt"

// Invalid: type switch without guard, re-asserting inside case clauses
func invalidTypeSwitch(x interface{}) {
	switch x.(type) {
	case int:
		fmt.Println(x.(int)) // MATCH /use the result of the type switch instead of asserting x.(int)/
	case string:
		fmt.Println(x.(string)) // MATCH /use the result of the type switch instead of asserting x.(string)/
	}
}

// Invalid: type switch with multiple assertions in one case
func invalidMultipleAssertions(x interface{}) {
	switch x.(type) {
	case int:
		a := x.(int) // MATCH /use the result of the type switch instead of asserting x.(int)/
		fmt.Println(a)
	}
}

// Valid: type switch with guard variable
func validTypeSwitchGuard(x interface{}) {
	switch v := x.(type) {
	case int:
		fmt.Println(v)
	case string:
		fmt.Println(v)
	}
}

// Valid: type switch without any type assertions in case body
func validNoAssertions(x interface{}) {
	switch x.(type) {
	case int:
		fmt.Println("int")
	case string:
		fmt.Println("string")
	}
}

// Valid: type assertion on a different variable
func validDifferentVar(x interface{}, y interface{}) {
	switch x.(type) {
	case int:
		fmt.Println(y.(int))
	}
}

// Valid: multi-type case clause (assertion would not match the guard)
func validMultiTypeCase(x interface{}) {
	switch x.(type) {
	case int, float64:
		fmt.Println("number")
	}
}
