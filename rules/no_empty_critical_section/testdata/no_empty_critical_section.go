package fixtures

import "sync"

var mu sync.Mutex

func emptyLockUnlock() {
	mu.Lock()    // MATCH /empty critical section, did you mean to defer the unlock?/
	mu.Unlock()
}

func deferredUnlock() {
	mu.Lock()
	defer mu.Unlock()
	// Work protected by the mutex
}

func workBetweenLockUnlock() {
	mu.Lock()
	_ = 42
	mu.Unlock()
}

var rw sync.RWMutex

func emptyRLockRUnlock() {
	rw.RLock()    // MATCH /empty critical section, did you mean to defer the unlock?/
	rw.RUnlock()
}

func validRLockDefer() {
	rw.RLock()
	defer rw.RUnlock()
	// Work protected by the read lock
}

type myStruct struct {
	mu sync.Mutex
}

func emptyStructMutex() {
	s := myStruct{}
	s.mu.Lock()    // MATCH /empty critical section, did you mean to defer the unlock?/
	s.mu.Unlock()
}

func validStructMutex() {
	s := myStruct{}
	s.mu.Lock()
	defer s.mu.Unlock()
	// Protected work
}

func differentReceivers() {
	var mu2 sync.Mutex
	mu.Lock()
	mu2.Unlock()
	_ = mu2
}
