package fixtures

import "time"

// Invalid: select with single case receiving from time.After and empty body
func waitElaborate() {
	select { // MATCH /use time.Sleep instead of select{case <-time.After(...)}/
	case <-time.After(1 * time.Second):
	}
}

// Invalid: select with single case receiving from time.After using a variable duration
func waitVariable() {
	d := 5 * time.Second
	select { // MATCH /use time.Sleep instead of select{case <-time.After(...)}/
	case <-time.After(d):
	}
}

// Valid: uses time.Sleep directly
func waitSimple() {
	time.Sleep(1 * time.Second)
}

// Valid: select with case that has a body
func waitWithBody() {
	x := 0
	select {
	case <-time.After(1 * time.Second):
		x++
	}
	_ = x
}

// Valid: select with multiple cases
func waitMultipleCases() {
	ch := make(chan int)
	select {
	case <-time.After(1 * time.Second):
	case <-ch:
	}
}

// Valid: select with case and default
func waitWithDefault() {
	select {
	case <-time.After(1 * time.Second):
	default:
	}
}

// Valid: select receiving from a regular channel
func waitRegularChannel() {
	ch := make(chan struct{})
	select {
	case <-ch:
	}
}

// Valid: select with assignment from time.After but with body
func waitAssignWithBody() {
	x := 0
	select {
	case _ = <-time.After(1 * time.Second):
		x++
	}
	_ = x
}
