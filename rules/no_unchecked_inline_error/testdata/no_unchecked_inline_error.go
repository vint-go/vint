package fixtures

func doSomething() (interface{}, error) { return nil, nil }
func use(v interface{})                 {}

// Invalid: error assigned but not checked in condition
func invalidUncheckedErr() {
	if val, err := doSomething(); val != nil { // MATCH /error variable 'err' assigned in if initialization but not checked in the condition/
		_ = err
		use(val)
	}
}

// Valid: error is checked in the condition
func validErrChecked() {
	if val, err := doSomething(); err == nil {
		use(val)
	}
}

// Valid: error variable used in condition (err != nil)
func validErrNotNil() {
	if _, err := doSomething(); err != nil {
		use(err)
	}
}

// Valid: separate assignment and check
func validSeparateCheck() {
	val, err := doSomething()
	if err != nil {
		return
	}
	use(val)
}

// Valid: no error variable in assignment
func validNoErrVar() {
	if val, ok := doSomething(); ok != nil {
		use(val)
	}
}

// Valid: no init clause
func validNoInit() {
	x, _ := doSomething()
	if x != nil {
		use(x)
	}
}

// Invalid: error variable assigned but condition checks something else
func invalidErrIgnored() {
	if _, err := doSomething(); true { // MATCH /error variable 'err' assigned in if initialization but not checked in the condition/
		_ = err
	}
}

// Valid: blank identifier for error (not flagged)
func validBlankErr() {
	if val, _ := doSomething(); val != nil {
		use(val)
	}
}
