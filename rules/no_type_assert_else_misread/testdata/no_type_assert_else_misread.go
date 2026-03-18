package fixtures

import "fmt"

// Case 1: separate assignment, ok condition, v used in else
func separateAssignElse(x interface{}) {
	v, ok := x.(string)
	if ok {
		fmt.Println("string:", v)
	} else {
		fmt.Println("not a string:", v) // MATCH /type assertion else branch reads v which is the zero value of the asserted type, not the original value/
	}
}

// Case 2: init form, v used in else
func initFormElse(x interface{}) {
	if v, ok := x.(string); ok {
		fmt.Println("string:", v)
	} else {
		fmt.Println("not a string:", v) // MATCH /type assertion else branch reads v which is the zero value of the asserted type, not the original value/
	}
}

// Case 3: negated condition, v used in if body (which is the else-equivalent)
func negatedCondition(x interface{}) {
	v, ok := x.(string)
	if !ok {
		fmt.Println("not a string:", v) // MATCH /type assertion else branch reads v which is the zero value of the asserted type, not the original value/
	} else {
		fmt.Println("string:", v)
	}
}

// Case 4: negated condition with init
func negatedConditionInit(x interface{}) {
	if v, ok := x.(string); !ok {
		fmt.Println("not a string:", v) // MATCH /type assertion else branch reads v which is the zero value of the asserted type, not the original value/
	} else {
		fmt.Println("string:", v)
	}
}

// Valid: v not used in else
func validNoElseUse(x interface{}) {
	v, ok := x.(string)
	if ok {
		fmt.Println("string:", v)
	} else {
		fmt.Println("not a string:", x) // uses x, not v
	}
}

// Valid: no else branch
func validNoElse(x interface{}) {
	v, ok := x.(string)
	if ok {
		fmt.Println("string:", v)
	}
	_ = ok
}

// Valid: single-value type assertion (no comma-ok)
func validSingleValue(x interface{}) {
	v := x.(string)
	fmt.Println("string:", v)
}

// Valid: ok is blank identifier
func validBlankOk(x interface{}) {
	v, _ := x.(string)
	fmt.Println("string:", v)
}

// Valid: v is blank identifier
func validBlankV(x interface{}) {
	_, ok := x.(string)
	if ok {
		fmt.Println("is string")
	} else {
		fmt.Println("not string")
	}
}

// Valid: condition is not the ok variable
func validDifferentCondition(x interface{}) {
	v, ok := x.(string)
	other := true
	if other {
		fmt.Println("other:", v)
	} else {
		fmt.Println("not other:", v)
	}
	_ = ok
}

// Valid: v used in the success branch, not the else branch
func validVInSuccessBranch(x interface{}) {
	v, ok := x.(string)
	if ok {
		fmt.Println("string:", v)
	} else {
		fmt.Println("not a string")
	}
}

// Case 5: multiple references in else
func multipleRefsInElse(x interface{}) {
	v, ok := x.(string)
	if ok {
		fmt.Println("string:", v)
	} else {
		fmt.Println("not a string:", v) // MATCH /type assertion else branch reads v which is the zero value of the asserted type, not the original value/
		fmt.Println("also:", v)         // MATCH /type assertion else branch reads v which is the zero value of the asserted type, not the original value/
	}
}
