package fixtures

import "errors"

// Invalid: err is shadowed by the inner := declaration in if-init.
func shadowedErrInIf() (err error) {
	x, err := doFirst()
	if err != nil {
		return err
	}
	if y, err := doSecond(x); err != nil { // MATCH /variable err shadows variable from outer scope/
		_ = y
		return err
	}
	return nil
}

// Invalid: variable shadowed in nested block.
func shadowedInBlock() {
	x := 1
	{
		x := 2 // MATCH /variable x shadows variable from outer scope/
		_ = x
	}
	_ = x
}

// Invalid: variable shadowed in for loop.
func shadowedInForLoop() {
	err := errors.New("outer")
	for i := 0; i < 3; i++ {
		err := errors.New("inner") // MATCH /variable err shadows variable from outer scope/
		_ = err
	}
	_ = err
}

// Invalid: variable shadowed in function literal.
func shadowedInFuncLit() {
	x := 1
	f := func() {
		x := 2 // MATCH /variable x shadows variable from outer scope/
		_ = x
	}
	f()
	_ = x
}

// Valid: use = instead of := to avoid shadowing.
func noShadowWithAssign() (err error) {
	x, err := doFirst()
	if err != nil {
		return err
	}
	var y int
	y, err = doSecond(x)
	if err != nil {
		return err
	}
	_ = y
	return nil
}

// Valid: different variable names.
func differentNames() {
	x := 1
	{
		y := 2
		_ = y
	}
	_ = x
}

// Valid: blank identifier is not flagged.
func blankIdentifier() {
	_ = 1
	{
		_ = 2
	}
}

// Valid: parameter and local variable have different names.
func paramNoShadow(a int) {
	b := a + 1
	_ = b
}

// helper stubs to make the file parse
func doFirst() (int, error)      { return 0, nil }
func doSecond(x int) (int, error) { return x, nil }
