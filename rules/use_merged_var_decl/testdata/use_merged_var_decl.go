package fixtures

// Invalid: var declaration followed by assignment
func bad1() {
	var x int // MATCH /should merge variable declaration with assignment on next line to x := .../
	x = 42
	_ = x
}

// Invalid: var declaration with type followed by assignment
func bad2() {
	var s string // MATCH /should merge variable declaration with assignment on next line to s := .../
	s = "hello"
	_ = s
}

// Invalid: var declaration of slice type
func bad3() {
	var items []int // MATCH /should merge variable declaration with assignment on next line to items := .../
	items = make([]int, 10)
	_ = items
}

// Valid: short variable declaration already used
func good1() {
	x := 42
	_ = x
}

// Valid: var with initial value
func good2() {
	var x int = 42
	_ = x
}

// Valid: var declaration without immediate assignment
func good3() {
	var x int
	println("something")
	x = 42
	_ = x
}

// Valid: assignment to different variable
func good4() {
	var x int
	y := 10
	x = y
	_ = x
}

// Valid: var declaration with multiple names
func good5() {
	var x, y int
	x = 1
	y = 2
	_ = x
	_ = y
}

// Valid: multi-value assignment
func good6() {
	var x int
	_ = x
}

// Valid: assignment uses +=
func good7() {
	var x int
	x += 1
	_ = x
}

// Valid: var block with multiple specs
func good8() {
	var (
		x int
		y int
	)
	x = 1
	_ = x
	_ = y
}

// Valid: variable is conditionally overridden (multiple assignments)
func good9() {
	var gormLogger int
	gormLogger = 1
	if true {
		gormLogger = 2
	}
	_ = gormLogger
}

// Valid: variable is reassigned inside a closure (multiple assignments)
func good10() {
	var err error
	err = doSomething(func() {
		err = nil
	})
	_ = err
}

// Valid: self-referential RHS
func good11() {
	var x int
	x = process(x)
	_ = x
}

// Valid: self-referential RHS nested in expression
func good12() {
	var s string
	s = s + "suffix"
	_ = s
}

// Valid: multiple assignments in different branches
func good13() {
	var x int
	x = 1
	for i := 0; i < 10; i++ {
		x = i
	}
	_ = x
}

func doSomething(f func()) error { f(); return nil }
func process(x int) int          { return x }
