package fixtures

func weakGuardNeqAnd(xs []int) {
	if xs != nil && xs[0] > 0 { // MATCH /nil check for 'xs' is not enough, use len(xs) != 0/
		println("first element is positive")
	}
}

func weakGuardEqlOr(xs []int) {
	if xs == nil || xs[0] > 0 { // MATCH /nil check for 'xs' is not enough, use len(xs) != 0/
		println("nil or first element is positive")
	}
}

func weakGuardNeqAndNilExpr(xs []*int) {
	if xs != nil && xs[0] != nil { // MATCH /nil check for 'xs' is not enough, use len(xs) != 0/
		println("has non-nil first element")
	}
}

func weakGuardReversedNil(xs []int) {
	if nil != xs && xs[0] > 0 { // MATCH /nil check for 'xs' is not enough, use len(xs) != 0/
		println("reversed nil check")
	}
}

func weakGuardReversedNilEql(xs []int) {
	if nil == xs || xs[0] > 0 { // MATCH /nil check for 'xs' is not enough, use len(xs) != 0/
		println("reversed nil eq check")
	}
}

// Valid: uses len instead of nil check
func goodLenGuard(xs []int) {
	if len(xs) != 0 && xs[0] > 0 {
		println("properly checks length before indexing")
	}
}

// Valid: uses len > 0
func goodLenGtGuard(xs []int) {
	if len(xs) > 0 && xs[0] > 0 {
		println("properly checks length before indexing")
	}
}

// Valid: nil check without indexing
func goodNilCheckOnly(xs []int) {
	if xs != nil {
		println("just a nil check")
	}
}

// Valid: indexing without nil check
func goodNoNilCheck(xs []int) {
	if xs[0] > 0 {
		println("no nil check")
	}
}

// Valid: different variable checked vs indexed
func goodDifferentVar(xs []int, ys []int) {
	if xs != nil && ys[0] > 0 {
		println("different variable")
	}
}

// Valid: equality check with && does not match the pattern
func goodEqlAnd(xs []int) {
	if xs == nil && len(xs) == 0 {
		println("equality with and")
	}
}

// Valid: inequality check with || does not match the pattern
func goodNeqOr(xs []int) {
	if xs != nil || len(xs) > 0 {
		println("inequality with or")
	}
}
