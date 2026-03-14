package fixtures

func longFunction() { // MATCH /function longFunction has too many lines (7 > 5)/
	a := 1
	b := 2
	c := 3
	d := 4
	e := 5
	f := 6
	println(a, b, c, d, e, f)
}

func shortFunction() {
	a := 1
	b := 2
	println(a, b)
}

func exactlyAtLimit() {
	a := 1
	b := 2
	c := 3
	d := 4
	println(a, b, c, d)
}

func emptyFunction() {}

func oneLineFunction() {
	println("hello")
}

func justOverLimit() { // MATCH /function justOverLimit has too many lines (6 > 5)/
	a := 1
	b := 2
	c := 3
	d := 4
	e := 5
	println(a, b, c, d, e)
}
