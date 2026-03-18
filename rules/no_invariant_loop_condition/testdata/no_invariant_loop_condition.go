package fixtures

func process(items []string) {}

// Invalid: condition variable i is never modified in the loop body
func invalidInvariantLoop() {
	i := 0
	for i < 10 { // MATCH /loop condition variable is never modified in the loop body/
		process(nil)
	}
}

// Invalid: condition variable with len call, variable never modified
func invalidInvariantLoopLen(items []string) {
	i := 0
	for i < len(items) { // MATCH /loop condition variable is never modified in the loop body/
		process(items)
	}
}

// Valid: condition variable is modified via post statement
func validPostIncrement() {
	for i := 0; i < 10; i++ {
		process(nil)
	}
}

// Valid: condition variable is modified in body
func validBodyIncrement() {
	i := 0
	for i < 10 {
		i++
	}
}

// Valid: condition variable is modified by assignment in body
func validBodyAssignment(items []string) {
	i := 0
	for i < len(items) {
		i = i + 1
	}
}

// Valid: infinite loop (no condition)
func validInfiniteLoop() {
	for {
		break
	}
}

// Valid: for range loop
func validRangeLoop(items []string) {
	for _, item := range items {
		process([]string{item})
	}
}
