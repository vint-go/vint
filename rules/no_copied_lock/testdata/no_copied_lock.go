package fixtures

import "sync"

type MyStruct struct {
	mu sync.Mutex
}

type SafeStruct struct {
	name string
}

// Bad: mutex is copied when passed by value
func badPassMutexByValue(mu sync.Mutex) { // MATCH /mu passes lock by value: sync.Mutex/
	_ = &mu
}

// Good: mutex is passed by pointer
func goodPassMutexByPointer(mu *sync.Mutex) {
	mu.Lock()
	defer mu.Unlock()
}

// Bad: struct containing mutex is passed by value
func badPassStructByValue(s MyStruct) { // MATCH /s passes lock by value: MyStruct.mu (sync.Mutex)/
	_ = &s
}

// Good: struct containing mutex passed by pointer
func goodPassStructByPointer(s *MyStruct) {
	_ = s
}

// Bad: assignment copies mutex inside struct
func badAssignmentCopy() {
	var a MyStruct
	b := a // MATCH /assignment copies lock value: MyStruct.mu (sync.Mutex)/
	_ = &b
}

// Good: pointer assignment
func goodPointerAssignment() {
	a := &MyStruct{}
	b := a // copies the pointer, not the mutex
	_ = b
}

// Good: composite literal creates new value
func goodCompositeLiteral() {
	a := MyStruct{}
	_ = &a
}

// Bad: return copies lock value
func badReturnCopy() MyStruct {
	var s MyStruct
	return s // MATCH /return copies lock value: MyStruct.mu (sync.Mutex)/
}

// Good: return pointer
func goodReturnPointer() *MyStruct {
	s := &MyStruct{}
	return s
}

// Bad: RWMutex passed by value
func badRWMutexByValue(mu sync.RWMutex) { // MATCH /mu passes lock by value: sync.RWMutex/
	_ = &mu
}

// Bad: WaitGroup passed by value
func badWaitGroupByValue(wg sync.WaitGroup) { // MATCH /wg passes lock by value: sync.WaitGroup/
	_ = &wg
}

// Bad: range copies lock value
func badRangeCopy() {
	structs := []MyStruct{{}, {}}
	for _, s := range structs { // MATCH /range var copies lock value: MyStruct.mu (sync.Mutex)/
		_ = &s
	}
}

// Good: range with index only
func goodRangeIndex() {
	structs := []MyStruct{{}, {}}
	for i := range structs {
		_ = &structs[i]
	}
}

// Good: range over pointers
func goodRangePointers() {
	structs := []*MyStruct{{}, {}}
	for _, s := range structs {
		_ = s
	}
}

// Good: plain struct without lock
func goodPlainStruct(s SafeStruct) {
	_ = s
}
