package fixtures

func process(item int) {}

// Invalid: for{} with if+break at start of body
func badInfiniteLoop() {
	items := []int{1, 2, 3}
	i := 0
	for {
		if i >= len(items) { // MATCH /lift if+break into loop condition/
			break
		}
		process(items[i])
		i++
	}
}

// Invalid: for{} with simple boolean condition
func badBoolCondition() {
	done := false
	for {
		if done { // MATCH /lift if+break into loop condition/
			break
		}
		done = true
	}
}

// Invalid: for{} with only if+break (no other statements after)
func badOnlyIfBreak() {
	x := 0
	for {
		if x > 10 { // MATCH /lift if+break into loop condition/
			break
		}
		x++
	}
}

// Valid: for loop already has a condition
func goodHasCondition() {
	items := []int{1, 2, 3}
	for i := 0; i < len(items); i++ {
		process(items[i])
	}
}

// Valid: if+break is not the first statement
func goodNotFirst() {
	i := 0
	for {
		process(i)
		if i >= 10 {
			break
		}
		i++
	}
}

// Valid: if body has more than just break
func goodIfBodyNotJustBreak() {
	i := 0
	for {
		if i >= 10 {
			process(i)
			break
		}
		i++
	}
}

// Valid: if has else clause
func goodIfHasElse() {
	i := 0
	for {
		if i >= 10 {
			break
		} else {
			process(i)
		}
		i++
	}
}

// Valid: break has a label
func goodBreakWithLabel() {
	i := 0
outer:
	for {
		if i >= 10 {
			break outer
		}
		i++
	}
}

// Valid: if has init statement
func goodIfHasInit() {
	items := []int{1, 2, 3}
	i := 0
	for {
		if j := i; j >= len(items) {
			break
		}
		process(items[i])
		i++
	}
}

// Valid: range loop (not a for statement with no condition)
func goodRangeLoop() {
	items := []int{1, 2, 3}
	for _, item := range items {
		process(item)
	}
}

// Valid: for loop with empty body
func goodEmptyBody() {
	for {
		break
	}
}
