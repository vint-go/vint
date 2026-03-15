package fixtures

import "sync"

func doubleLock() {
	var mu sync.Mutex
	mu.Lock()
	mu.Lock() // MATCH /suspicious lock sequence: Lock called twice without Unlock/
}

func doubleUnlock() {
	var mu sync.Mutex
	mu.Unlock()
	mu.Unlock() // MATCH /suspicious lock sequence: Unlock called twice without Lock/
}

func emptyLockUnlock() {
	var mu sync.Mutex
	mu.Lock()
	mu.Unlock() // MATCH /suspicious lock sequence: Lock immediately followed by Unlock with no work in critical section/
}

func doubleRLock() {
	var rw sync.RWMutex
	rw.RLock()
	rw.RLock() // MATCH /suspicious lock sequence: RLock called twice without RUnlock/
}

func doubleRUnlock() {
	var rw sync.RWMutex
	rw.RUnlock()
	rw.RUnlock() // MATCH /suspicious lock sequence: RUnlock called twice without RLock/
}

func emptyRLockRUnlock() {
	var rw sync.RWMutex
	rw.RLock()
	rw.RUnlock() // MATCH /suspicious lock sequence: RLock immediately followed by RUnlock with no work in critical section/
}

// Valid patterns below

func validLockDefer() {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	// critical section
}

func validLockWithWork() {
	var mu sync.Mutex
	mu.Lock()
	doWork()
	mu.Unlock()
}

func validRLockDefer() {
	var rw sync.RWMutex
	rw.RLock()
	defer rw.RUnlock()
	// critical section
}

func validRLockWithWork() {
	var rw sync.RWMutex
	rw.RLock()
	doWork()
	rw.RUnlock()
}

func validDifferentMutexes() {
	var mu1, mu2 sync.Mutex
	mu1.Lock()
	mu2.Lock()
	doWork()
	mu2.Unlock()
	mu1.Unlock()
}

func doWork() {}
