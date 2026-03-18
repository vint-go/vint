package fixtures

// Invalid: assignment in if init where short var decl would be better.
func bad1() error {
	var err error
	if err = doSomething(); err != nil { // MATCH /use short variable declaration (:=) instead of assignment (=) in if init/
		return err
	}
	return nil
}

// Invalid: same pattern with a different variable name.
func bad2() error {
	var e error
	if e = doAnother(); e != nil { // MATCH /use short variable declaration (:=) instead of assignment (=) in if init/
		return e
	}
	return nil
}

// Valid: already using short variable declaration.
func good1() error {
	if err := doSomething(); err != nil {
		return err
	}
	return nil
}

// Valid: no init statement.
func good2() error {
	err := doSomething()
	if err != nil {
		return err
	}
	return nil
}

// Valid: condition is not != nil.
func good3() error {
	var err error
	if err = doSomething(); err == nil {
		return err
	}
	return nil
}

// Valid: body does not return the variable.
func good4() error {
	var err error
	if err = doSomething(); err != nil {
		return otherErr()
	}
	_ = err
	return nil
}

// Valid: multiple assignments in init.
func good5() (int, error) {
	var n int
	var err error
	_ = n
	_ = err
	return 0, nil
}

// helper stubs to make the file parse
func doSomething() error { return nil }
func doAnother() error   { return nil }
func otherErr() error    { return nil }
