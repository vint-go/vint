package fixtures

import "fmt"

// --- Invalid: redundant return at end of void function ---

func redundantReturn() {
	fmt.Println("done")
	return // MATCH /redundant return statement/
}

func redundantReturnEmpty() {
	return // MATCH /redundant return statement/
}

func redundantReturnMultiStmt() {
	x := 1
	_ = x
	fmt.Println("hello")
	return // MATCH /redundant return statement/
}

// --- Invalid: redundant break at end of case clause ---

func redundantBreakSwitch() {
	x := 1
	switch x {
	case 1:
		fmt.Println("one")
		break // MATCH /redundant break statement/
	case 2:
		fmt.Println("two")
	}
}

func redundantBreakDefault() {
	x := 1
	switch x {
	case 1:
		fmt.Println("one")
	default:
		fmt.Println("default")
		break // MATCH /redundant break statement/
	}
}

func redundantBreakTypeSwitch() {
	var x interface{} = 1
	switch x.(type) {
	case int:
		fmt.Println("int")
		break // MATCH /redundant break statement/
	case string:
		fmt.Println("string")
	}
}

// --- Invalid: redundant return in function literal ---

var _ = func() {
	fmt.Println("literal")
	return // MATCH /redundant return statement/
}

// --- Valid: function with return value ---

func validReturnValue() int {
	return 42
}

// --- Valid: return not at end ---

func validReturnMiddle() {
	if true {
		return
	}
	fmt.Println("still going")
}

// --- Valid: no return at all ---

func validNoReturn() {
	fmt.Println("done")
}

// --- Valid: break with label ---

func validLabeledBreak() {
outer:
	for i := 0; i < 10; i++ {
		switch i {
		case 5:
			break outer
		}
	}
}

// --- Valid: case with no break ---

func validNoCaseBreak() {
	x := 1
	switch x {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	}
}

// --- Valid: named results with bare return is still a return of values ---

func validNamedReturn() (result int) {
	result = 42
	return
}

// --- Valid: break-only case clause (explicit no-op, not redundant) ---

func validBreakOnlyCase() {
	x := 1
	switch x {
	case 1:
		break
	case 2:
		fmt.Println("two")
	}
}

// --- Valid: select statement break (not checked) ---

func validSelectBreak() {
	ch := make(chan int, 1)
	ch <- 1
	select {
	case v := <-ch:
		fmt.Println(v)
		break
	}
}

// --- Valid: break not at end of case ---

func validBreakMiddle() {
	x := 1
	switch x {
	case 1:
		if x > 0 {
			break
		}
		fmt.Println("one")
	}
}
