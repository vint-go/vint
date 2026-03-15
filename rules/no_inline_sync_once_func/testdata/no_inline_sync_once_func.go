package fixtures

import "sync"

func initialize() {}

// Invalid: calling sync.OnceFunc inline defeats its purpose
func bad() {
	sync.OnceFunc(func() { // MATCH /sync.OnceFunc result is called inline and should be stored in a variable instead/
		initialize()
	})()
}

// Valid: storing the result in a variable
var initOnce = sync.OnceFunc(func() {
	initialize()
})

func good() {
	initOnce() // safe to call multiple times
}

// Valid: assigning to a variable
func goodLocal() {
	f := sync.OnceFunc(func() {
		initialize()
	})
	f()
}
