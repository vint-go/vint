package fixtures

func process(items []string) {
	for _, item := range items {
		switch item {
		case "stop":
			// Only breaks out of switch, not the for loop
			break // MATCH /ineffective break statement. Did you mean to break out of the outer loop?/
		}
	}
}

func processSelect(ch chan int, done chan bool) {
	for {
		select {
		case <-ch:
			break // MATCH /ineffective break statement. Did you mean to break out of the outer loop?/
		case <-done:
			return
		}
	}
}

func processForStmt(items []string) {
	for i := 0; i < len(items); i++ {
		switch items[i] {
		case "stop":
			break // MATCH /ineffective break statement. Did you mean to break out of the outer loop?/
		case "continue":
			doSomething()
		}
	}
}

func processTypeSwitch(items []interface{}) {
	for _, item := range items {
		switch item.(type) {
		case string:
			break // MATCH /ineffective break statement. Did you mean to break out of the outer loop?/
		case int:
			doSomething()
		}
	}
}

// Valid: labeled break
func processLabeled(items []string) {
loop:
	for _, item := range items {
		switch item {
		case "stop":
			// Breaks out of the for loop
			break loop
		}
	}
}

// Valid: no enclosing loop
func processNoLoop(x int) {
	switch x {
	case 1:
		break
	case 2:
		doSomething()
	}
}

// Valid: no enclosing loop select
func processNoLoopSelect(ch chan int, done chan bool) {
	select {
	case <-ch:
		break
	case <-done:
		return
	}
}

// Invalid: break in nested loop's switch is also ineffective
func processNestedLoop(items []string) {
	for _, item := range items {
		for _, c := range item {
			switch {
			case c == 'a':
				break // MATCH /ineffective break statement. Did you mean to break out of the outer loop?/
			}
		}
	}
}

// Valid: no break
func processNoBreak(items []string) {
	for _, item := range items {
		switch item {
		case "stop":
			return
		}
	}
}

// Invalid: break inside if within switch case in a loop
func processBreakInIf(items []string) {
	for _, item := range items {
		switch item {
		case "stop":
			if len(item) > 0 {
				break // MATCH /ineffective break statement. Did you mean to break out of the outer loop?/
			}
		}
	}
}

func doSomething() {}
