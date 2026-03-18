package fixtures

import "fmt"

func redundantWithLiterals() {
	var x int = 42 // MATCH /redundant type int in variable declaration/
	_ = x

	var s string = "hello" // MATCH /redundant type string in variable declaration/
	_ = s

	var b bool = true // MATCH /redundant type bool in variable declaration/
	_ = b

	var f float64 = 1.0 // MATCH /redundant type float64 in variable declaration/
	_ = f

	var c complex128 = 1i // MATCH /redundant type complex128 in variable declaration/
	_ = c
}

func redundantWithTypedExpr() {
	var x int = int(42) // MATCH /redundant type int in variable declaration/
	_ = x

	var u uint = uint(42) // MATCH /redundant type uint in variable declaration/
	_ = u
}

func redundantWithFuncCall() {
	var s string = fmt.Sprintf("hello %s", "world") // MATCH /redundant type string in variable declaration/
	_ = s
}

func valid() {
	// No explicit type -- ok
	var x = 42
	_ = x

	// Short variable declaration -- ok
	y := 42
	_ = y

	// Explicit type with no value -- ok
	var z int
	_ = z

	// Type differs from default RHS type -- not redundant
	var f float64 = 1
	_ = f

	// Interface type with concrete value -- not redundant
	var i interface{} = 42
	_ = i

	fmt.Println("ok")
}
