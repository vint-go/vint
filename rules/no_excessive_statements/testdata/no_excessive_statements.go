package fixtures

func tooManyStatements() { // MATCH /function tooManyStatements has too many statements (6 > 5)/
	a := 1
	b := 2
	c := 3
	d := 4
	e := 5
	println(a, b, c, d, e)
}

func justRight() {
	a := 1
	b := 2
	c := 3
	d := 4
	println(a, b, c, d)
}

func emptyFunction() {}

func withControlFlow() { // MATCH /function withControlFlow has too many statements (6 > 5)/
	a := 1
	if a > 0 {
		b := 2
		println(b)
	}
	c := 3
	println(a, c)
}

func withForLoop() { // MATCH /function withForLoop has too many statements (7 > 5)/
	sum := 0
	for i := 0; i < 10; i++ {
		sum += i
		println(sum)
	}
	a := 1
	b := 2
	println(a, b)
}

func withGoStatement() { // MATCH /function withGoStatement has too many statements (8 > 5)/
	a := 1
	go func() {
		b := 2
		c := 3
		println(b, c)
	}()
	d := 4
	e := 5
	println(a, d, e)
}

func withDeferStatement() { // MATCH /function withDeferStatement has too many statements (7 > 5)/
	a := 1
	defer func() {
		b := 2
		println(b)
	}()
	c := 3
	d := 4
	println(a, c, d)
}

func shortWithControlFlow() {
	a := 1
	if a > 0 {
		println(a)
	}
}

func withSwitch() { // MATCH /function withSwitch has too many statements (7 > 5)/
	a := 1
	switch a {
	case 1:
		println("one")
		println("1")
	case 2:
		println("two")
		println("2")
	}
	println(a)
}
