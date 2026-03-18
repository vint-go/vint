package fixtures

func setX(x *int) error {
	*x = 42
	return nil
}

func setY(y *int) error {
	*y = 99
	return nil
}

// Invalid: x may be modified by setX via &x
func badReturnAddrOf() (int, error) {
	x := 10
	return x, setX(&x) // MATCH /return value may depend on evaluation order: variable x may be modified by a function call in another return value/
}

// Invalid: y may be modified by setY via &y
func badReturnAddrOfY() (int, error) {
	y := 20
	return y, setY(&y) // MATCH /return value may depend on evaluation order: variable y may be modified by a function call in another return value/
}

// Valid: value is stored before return
func goodReturnSeparate() (int, error) {
	x := 10
	err := setX(&x)
	return x, err
}

// Valid: no address-of in call arguments
func goodReturnNoAddr() (int, error) {
	x := 10
	return x, nil
}

// Valid: single return value
func goodSingleReturn() int {
	x := 10
	return x
}

// Valid: different variables, no dependency
func goodDifferentVars() (int, error) {
	x := 10
	y := 20
	return x, setY(&y)
}
