package fixtures

func doSomething() {}

// Invalid: nil on the left side of comparison
func yodaNilCheck(err error) error {
	if nil != err { // MATCH /yoda condition: rewrite nil != err as err != nil/
		return err
	}
	return nil
}

// Invalid: numeric literal on the left side
func yodaNumericCheck(x int) {
	if 42 == x { // MATCH /yoda condition: rewrite 42 == x as x == 42/
		doSomething()
	}
}

// Invalid: string literal on the left side
func yodaStringCheck(s string) {
	if "hello" == s { // MATCH /yoda condition: rewrite "hello" == s as s == "hello"/
		doSomething()
	}
}

// Invalid: boolean literal on the left side
func yodaBoolCheck(b bool) {
	if true == b { // MATCH /yoda condition: rewrite true == b as b == true/
		doSomething()
	}
}

// Invalid: negative number on the left side
func yodaNegativeCheck(x int) {
	if -1 == x { // MATCH /yoda condition: rewrite -1 == x as x == -1/
		doSomething()
	}
}

// Invalid: less-than with literal on left
func yodaLessThan(x int) {
	if 10 < x { // MATCH /yoda condition: rewrite 10 < x as x > 10/
		doSomething()
	}
}

// Invalid: greater-than with literal on left
func yodaGreaterThan(x int) {
	if 10 > x { // MATCH /yoda condition: rewrite 10 > x as x < 10/
		doSomething()
	}
}

// Valid: variable on the left side
func normalNilCheck(err error) error {
	if err != nil {
		return err
	}
	return nil
}

// Valid: variable on the left side
func normalNumericCheck(x int) {
	if x == 42 {
		doSomething()
	}
}

// Valid: two variables
func twoVarCheck(x, y int) {
	if x == y {
		doSomething()
	}
}

// Valid: two literals (both constant)
func twoLiteralCheck() {
	if 1 == 1 {
		doSomething()
	}
}

// Valid: function call on the left side
func funcCallCheck(x int) {
	if len("test") == x {
		doSomething()
	}
}
