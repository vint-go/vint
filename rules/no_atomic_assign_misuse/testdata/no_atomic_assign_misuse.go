package fixtures

import "sync/atomic"

func atomicAssignMisuseInvalid() {
	var x int64
	// Wrong: result of AddInt64 is assigned back to x
	x = atomic.AddInt64(&x, 1) // MATCH /direct assignment to atomic value/
}

func atomicAssignMisuseInvalidInt32() {
	var x int32
	x = atomic.AddInt32(&x, 1) // MATCH /direct assignment to atomic value/
}

func atomicAssignMisuseInvalidUint32() {
	var x uint32
	x = atomic.AddUint32(&x, 1) // MATCH /direct assignment to atomic value/
}

func atomicAssignMisuseInvalidUint64() {
	var x uint64
	x = atomic.AddUint64(&x, 1) // MATCH /direct assignment to atomic value/
}

func atomicAssignMisuseInvalidUintptr() {
	var x uintptr
	x = atomic.AddUintptr(&x, 1) // MATCH /direct assignment to atomic value/
}

func atomicAssignMisuseValidNewVar() {
	var x int64
	// Correct: assign to a new variable
	newVal := atomic.AddInt64(&x, 1)
	_ = newVal
}

func atomicAssignMisuseValidSideEffect() {
	var x int64
	// Correct: call for side effect only
	atomic.AddInt64(&x, 1)
}

func atomicAssignMisuseValidDifferentVar() {
	var x int64
	var y int64
	// Correct: assigned to a different variable
	y = atomic.AddInt64(&x, 1)
	_ = y
}
