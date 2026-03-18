package fixtures

import "os"

// Invalid: defer Close before checking error
func noDeferCloseBeforeErrCheckBad() {
	f, err := os.Open("file.txt")
	defer f.Close() // MATCH /possible nil dereference: Close deferred before checking error from assignment/
	if err != nil {
		return
	}
}

// Valid: error checked before defer Close
func noDeferCloseBeforeErrCheckGood() {
	f, err := os.Open("file.txt")
	if err != nil {
		return
	}
	defer f.Close()
}

// Valid: no error variable in assignment
func noDeferCloseBeforeErrCheckNoErr() {
	f, _ := os.Open("file.txt")
	_ = f
}

// Invalid: defer with wrapped Close in function literal
func noDeferCloseBeforeErrCheckWrapped() {
	f, err := os.Open("file.txt")
	defer func() { f.Close() }() // MATCH /possible nil dereference: Close deferred before checking error from assignment/
	if err != nil {
		return
	}
}

// Valid: error checked before defer in function literal
func noDeferCloseBeforeErrCheckWrappedGood() {
	f, err := os.Open("file.txt")
	if err != nil {
		return
	}
	defer func() { f.Close() }()
}

// Valid: defer on a different variable
func noDeferCloseBeforeErrCheckDifferentVar() {
	g, _ := os.Open("other.txt")
	f, err := os.Open("file.txt")
	defer g.Close()
	if err != nil {
		return
	}
	_ = f
}

// Invalid: inside nested function literal
func noDeferCloseBeforeErrCheckNestedFunc() {
	fn := func() {
		f, err := os.Open("file.txt")
		defer f.Close() // MATCH /possible nil dereference: Close deferred before checking error from assignment/
		if err != nil {
			return
		}
	}
	fn()
}

// Valid: single return value, no err
func noDeferCloseBeforeErrCheckSingleReturn() {
	f := os.NewFile(0, "test")
	defer f.Close()
}
