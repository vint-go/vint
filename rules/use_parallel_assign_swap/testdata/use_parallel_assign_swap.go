package fixtures

func swapBad() {
	a := 1
	b := 2
	tmp := a // MATCH /use parallel assignment to swap values: a, b = b, a/
	a = b
	b = tmp
}

func swapGood() {
	a := 1
	b := 2
	a, b = b, a
}

func swapBadWithDifferentNames() {
	x := 10
	y := 20
	temp := x // MATCH /use parallel assignment to swap values: x, y = y, x/
	x = y
	y = temp
}

func notASwap() {
	a := 1
	b := 2
	c := a
	a = b
	b = c + 1 // not a swap: RHS is not the temp variable
}

func notASwapDifferentAssignOp() {
	a := 1
	b := 2
	tmp := a
	a += b // not a plain assignment
	b = tmp
}

func multipleSwapsInBlock() {
	a := 1
	b := 2
	c := 3
	d := 4

	tmp := a // MATCH /use parallel assignment to swap values: a, b = b, a/
	a = b
	b = tmp

	tmp2 := c // MATCH /use parallel assignment to swap values: c, d = d, c/
	c = d
	d = tmp2
}

func swapFieldAccess() {
	type pair struct{ x, y int }
	p := pair{1, 2}
	tmp := p.x // MATCH /use parallel assignment to swap values: p.x, p.y = p.y, p.x/
	p.x = p.y
	p.y = tmp
}

func tooFewStatements() {
	a := 1
	_ = a
}
