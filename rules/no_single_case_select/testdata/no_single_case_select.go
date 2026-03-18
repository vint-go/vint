package fixtures

func singleCaseReceive() {
	ch := make(chan int)
	select { // MATCH /use plain channel send or receive instead of single-case select/
	case v := <-ch:
		_ = v
	}
}

func singleCaseSend() {
	ch := make(chan int)
	select { // MATCH /use plain channel send or receive instead of single-case select/
	case ch <- 1:
	}
}

func singleCaseReceiveDiscard() {
	ch := make(chan int)
	select { // MATCH /use plain channel send or receive instead of single-case select/
	case <-ch:
	}
}

// Valid: select with multiple cases
func multiCaseSelect() {
	ch1 := make(chan int)
	ch2 := make(chan int)
	select {
	case <-ch1:
	case <-ch2:
	}
}

// Valid: select with case and default
func caseWithDefault() {
	ch := make(chan int)
	select {
	case v := <-ch:
		_ = v
	default:
	}
}

// Valid: select with only default
func defaultOnly() {
	select {
	default:
	}
}

// Valid: plain channel receive
func plainReceive() {
	ch := make(chan int)
	v := <-ch
	_ = v
}

// Valid: plain channel send
func plainSend() {
	ch := make(chan int)
	ch <- 1
}
