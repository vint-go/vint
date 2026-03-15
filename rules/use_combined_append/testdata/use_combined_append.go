package fixtures

func consecutiveAppendInvalid() {
	var xs []int
	xs = append(xs, 1)
	xs = append(xs, 2) // MATCH /consecutive append to xs can be combined into a single call/
}

func consecutiveAppendThreeInvalid() {
	var result []int
	a, b, c := 1, 2, 3
	result = append(result, a)
	result = append(result, b) // MATCH /consecutive append to result can be combined into a single call/
	result = append(result, c)
}

func combinedAppendValid() {
	var xs []int
	xs = append(xs, 1, 2)
}

func combinedAppendMultipleValid() {
	var result []int
	a, b, c := 1, 2, 3
	result = append(result, a, b, c)
}

func appendWithBreakValid() {
	var xs []int
	xs = append(xs, 1)
	doSomething()
	xs = append(xs, 2)
}

func appendDifferentSlicesValid() {
	var xs, ys []int
	xs = append(xs, 1)
	ys = append(ys, 2)
}

func appendVariadicValid() {
	var xs []int
	extra := []int{1, 2, 3}
	xs = append(xs, 1)
	xs = append(xs, extra...)
}

func appendSingleArgValid() {
	var xs []int
	xs = append(xs, 1)
	_ = xs
}

func doSomething() {}
