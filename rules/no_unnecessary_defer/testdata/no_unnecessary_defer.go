package fixtures

import "sync"

var mu sync.Mutex

// Bad: defer immediately before return
func badDeferBeforeReturn() error {
	mu.Lock()
	defer mu.Unlock() // MATCH /unnecessary defer: right before return/
	return nil
}

// Bad: defer immediately before return with named return
func badDeferBeforeReturnNamed() (err error) {
	mu.Lock()
	defer mu.Unlock() // MATCH /unnecessary defer: right before return/
	return
}

// Bad: defer right before return in a nested if block
func badDeferInIf(cond bool) error {
	if cond {
		mu.Lock()
		defer mu.Unlock() // MATCH /unnecessary defer: right before return/
		return nil
	}
	return nil
}

// Good: defer with code after it
func goodDeferWithCodeAfter() error {
	mu.Lock()
	defer mu.Unlock()
	// ... do work ...
	doSomething()
	return nil
}

// Good: direct call instead of defer
func goodDirectCall() error {
	mu.Lock()
	mu.Unlock()
	return nil
}

// Good: defer at start of function with multiple statements after
func goodDeferAtStart() error {
	mu.Lock()
	defer mu.Unlock()
	x := doSomething()
	if x != nil {
		return x
	}
	return nil
}

// Good: only one statement (defer without return)
func goodOnlyDefer() {
	defer mu.Unlock()
}

// Good: return without preceding defer
func goodReturnOnly() error {
	return nil
}

func doSomething() error {
	return nil
}
