package fixtures

import "time"

func timerResetRetvalInvalid() {
	t := time.NewTimer(1 * time.Second)
	// Wrong: using Reset's return value in an if condition
	if !t.Reset(2 * time.Second) { // MATCH /the return value of (*time.Timer).Reset is unreliable and should not be used/
		<-t.C
	}
}

func timerResetRetvalAssign() {
	t := time.NewTimer(1 * time.Second)
	// Wrong: assigning Reset's return value
	wasActive := t.Reset(2 * time.Second) // MATCH /the return value of (*time.Timer).Reset is unreliable and should not be used/
	_ = wasActive
}

func timerResetRetvalValid() {
	t := time.NewTimer(1 * time.Second)
	if !t.Stop() {
		<-t.C
	}
	t.Reset(2 * time.Second)
}

func timerResetRetvalValidStandalone() {
	t := time.NewTimer(1 * time.Second)
	// OK: not using the return value
	t.Reset(3 * time.Second)
}
