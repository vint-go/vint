package fixtures

import "time"

// Invalid: time.Tick in a regular function leaks the ticker.
func process() {
	for range time.Tick(1 * time.Second) { // MATCH /calling time.Tick leaks the underlying ticker; use time.NewTicker and call Stop() when done/
		// do work
	}
}

// Invalid: time.Tick in a method.
func (s someStruct) doWork() {
	ch := time.Tick(500 * time.Millisecond) // MATCH /calling time.Tick leaks the underlying ticker; use time.NewTicker and call Stop() when done/
	_ = ch
}

// Invalid: time.Tick used in a variable assignment inside a helper.
func helper() {
	ticker := time.Tick(2 * time.Second) // MATCH /calling time.Tick leaks the underlying ticker; use time.NewTicker and call Stop() when done/
	_ = ticker
}

// Valid: time.Tick in main is acceptable.
func main() {
	for range time.Tick(1 * time.Second) {
		// long-lived function, acceptable
	}
}

// Valid: time.Tick in init is acceptable.
func init() {
	_ = time.Tick(1 * time.Second)
}

// Valid: using time.NewTicker with Stop.
func goodProcess() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		// do work
	}
}

type someStruct struct{}
