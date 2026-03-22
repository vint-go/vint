package fixtures

import "errors"

// In non-strict mode, only shadows where the outer variable is used
// after the shadowing scope should be reported.

// Invalid (non-strict): outer x is used after the shadowing block.
func shadowUsedAfter() {
	x := 1
	{
		x := 2 // MATCH /variable x shadows variable from outer scope/
		_ = x
	}
	_ = x // outer x used after
}

// Valid (non-strict): outer x is NOT used after the shadowing block.
func shadowNotUsedAfter() {
	x := 1
	_ = x
	{
		x := 2
		_ = x
	}
	// outer x is not used here
}

// Invalid (non-strict): outer err is used after the if scope.
func shadowErrUsedAfter() (err error) {
	x, err := doFirst()
	if err != nil {
		return err
	}
	if y, err := doSecond(x); err != nil { // MATCH /variable err shadows variable from outer scope/
		_ = y
		return err
	}
	return err // outer err used after the if scope
}

// Invalid (non-strict): outer err used after the for loop.
func shadowInForUsedAfter() {
	err := errors.New("outer")
	for i := 0; i < 3; i++ {
		err := errors.New("inner") // MATCH /variable err shadows variable from outer scope/
		_ = err
	}
	_ = err // used after for loop
}

// Valid (non-strict): outer err NOT used after the for loop.
func shadowInForNotUsedAfter() {
	err := errors.New("outer")
	_ = err
	for i := 0; i < 3; i++ {
		err := errors.New("inner")
		_ = err
	}
}

// Valid: different variable names (same in both modes).
func differentNamesNonStrict() {
	x := 1
	{
		y := 2
		_ = y
	}
	_ = x
}

// Valid: blank identifier is not flagged (same in both modes).
func blankIdentifierNonStrict() {
	_ = 1
	{
		_ = 2
	}
}

// helper stubs to make the file parse
func doFirst() (int, error)       { return 0, nil }
func doSecond(x int) (int, error) { return x, nil }
