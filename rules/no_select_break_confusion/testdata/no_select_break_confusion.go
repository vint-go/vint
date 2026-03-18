package fixtures

// --- Invalid examples ---

func unlabeledBreakInSelectInFor(done chan struct{}) {
	for {
		select {
		case <-done:
			break // MATCH /break inside select inside for loop only breaks the select, not the loop; use a labeled break to exit the loop/
		}
	}
}

func unlabeledBreakInSelectInForRange(ch chan int, items []int) {
	for range items {
		select {
		case v := <-ch:
			_ = v
			break // MATCH /break inside select inside for loop only breaks the select, not the loop; use a labeled break to exit the loop/
		}
	}
}

func multipleBreaksInSelect(done chan struct{}, data chan int) {
	for {
		select {
		case <-done:
			break // MATCH /break inside select inside for loop only breaks the select, not the loop; use a labeled break to exit the loop/
		case <-data:
			break // MATCH /break inside select inside for loop only breaks the select, not the loop; use a labeled break to exit the loop/
		}
	}
}

func breakInIfInsideSelect(done chan struct{}) {
	for {
		select {
		case <-done:
			if true {
				break // MATCH /break inside select inside for loop only breaks the select, not the loop; use a labeled break to exit the loop/
			}
		}
	}
}

func breakInForCondition(done chan struct{}) {
	for i := 0; i < 10; i++ {
		select {
		case <-done:
			break // MATCH /break inside select inside for loop only breaks the select, not the loop; use a labeled break to exit the loop/
		}
	}
}

// --- Valid examples ---

func labeledBreakInSelect(done chan struct{}) {
loop:
	for {
		select {
		case <-done:
			break loop
		}
	}
}

func breakInNestedForInsideSelect(done chan struct{}) {
	for {
		select {
		case <-done:
			for i := 0; i < 5; i++ {
				break // this breaks the inner for, not confusing
			}
		}
	}
}

func breakInSwitchInsideSelect(done chan struct{}, x int) {
	for {
		select {
		case <-done:
			switch x {
			case 1:
				break // this breaks the switch, not confusing
			}
		}
	}
}

func breakInSelectNotInFor(done chan struct{}) {
	select {
	case <-done:
		break // not inside a for loop, so not confusing
	}
}

func returnInSelectInFor(done chan struct{}) {
	for {
		select {
		case <-done:
			return // return is fine, not a break
		}
	}
}

func noBreakInSelect(done chan struct{}, data chan int) {
	results := make([]int, 0)
	for {
		select {
		case <-done:
			return
		case v := <-data:
			results = append(results, v)
		}
	}
}

func breakInFuncLitInsideSelect(done chan struct{}) {
	for {
		select {
		case <-done:
			func() {
				// break inside func literal is a compile error,
				// but for completeness: it wouldn't target the select
			}()
		}
	}
}
