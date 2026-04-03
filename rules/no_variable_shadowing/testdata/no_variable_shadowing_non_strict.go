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

// Valid (non-strict): if-init err := followed by another if-init err :=.
// The later "if err := ..." is a new declaration, NOT a use of the outer err.
// This was a false positive before types.Object-based usage checking.
func ifInitErrChain() {
	err := errors.New("outer")
	if err != nil {
		_ = err
	}
	if err := doThird(); err != nil { // shadows outer err
		_ = err
	}
	if err := doThird(); err != nil { // another new declaration, not a use of outer err
		_ = err
	}
	// outer err is NOT used after the shadowing if-blocks
}

// Valid (non-strict): err declared inside a closure shadows the outer err,
// but the outer err is only re-declared (not used) in subsequent statements.
func closureErrShadowWithLaterRedecl() {
	err := errors.New("setup")
	_ = err
	f := func() error {
		err := doThird() // shadows outer err
		return err
	}
	_ = f
	// This is a new err, not a use of the outer one.
	if err := doThird(); err != nil {
		_ = err
	}
}

// Valid (non-strict): variables inside a closure shadow outer variables,
// but the outer variables are only re-declared later, not used.
func closureVarShadowNoOuterUse() {
	ep := 1
	_ = ep
	f := func() {
		ep := 2 // shadows outer ep
		_ = ep
	}
	f()
	// ep1, ep2 are different names, and ep is not used after the closure
}

// Valid (non-strict): err inside a closure shadows outer err.
// Even though the outer err IS used after in the enclosing function,
// go vet's shadow resolves non-strict shadows within the function where
// the := occurs. Inside the closure, the outer err is not used after
// the shadow (the inner err takes over), so it is not reported.
func closureErrShadowWithOuterUse() (err error) {
	err = doThird()
	if err != nil {
		return
	}
	f := func() {
		err := doThird() // NOT reported: shadow is resolved within closure scope
		_ = err
	}
	f()
	return err // outer err used in enclosing function, but irrelevant to closure shadow
}

// Valid (non-strict): err in Transaction-like closure pattern.
// The outer err is the return of Transaction(); the inner err is a separate
// variable inside the closure. The outer err is not used after Transaction().
func transactionClosurePattern() {
	err := doTransaction(func() error {
		err := doThird() // shadows outer err
		return err
	})
	_ = err
}

// Invalid (non-strict): outer err IS used after the block that shadows it.
func blockShadowWithSubsequentUse() {
	err := doThird()
	{
		err := errors.New("inner") // MATCH /variable err shadows variable from outer scope/
		_ = err
	}
	if err != nil { // genuine use of outer err
		_ = err
	}
}

// helper stubs to make the file parse
func doFirst() (int, error)       { return 0, nil }
func doSecond(x int) (int, error) { return x, nil }
func doThird() error              { return nil }
func doTransaction(fn func() error) error {
	return fn()
}
