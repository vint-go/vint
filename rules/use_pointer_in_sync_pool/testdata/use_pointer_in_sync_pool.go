package fixtures

import "sync"

// Invalid: Put with a non-pointer value (slice is non-pointer for interface boxing).
func badPutSlice() {
	var pool sync.Pool
	b := make([]byte, 1024)
	pool.Put(b) // MATCH /non-pointer value stored in sync.Pool causes allocation; use a pointer instead/
}

// Invalid: Put with a struct value.
func badPutStruct() {
	type data struct{ x int }
	var pool sync.Pool
	pool.Put(data{x: 1}) // MATCH /non-pointer value stored in sync.Pool causes allocation; use a pointer instead/
}

// Invalid: Put with an int value.
func badPutInt() {
	var pool sync.Pool
	pool.Put(42) // MATCH /non-pointer value stored in sync.Pool causes allocation; use a pointer instead/
}

// Invalid: Put with a string value.
func badPutString() {
	var pool sync.Pool
	pool.Put("hello") // MATCH /non-pointer value stored in sync.Pool causes allocation; use a pointer instead/
}

// Valid: Put with a pointer value.
func goodPutPointer() {
	var pool sync.Pool
	b := make([]byte, 1024)
	pool.Put(&b)
}

// Valid: Put with a pointer to struct.
func goodPutStructPointer() {
	type data struct{ x int }
	var pool sync.Pool
	pool.Put(&data{x: 1})
}

// Valid: Put with an interface value.
func goodPutInterface() {
	var pool sync.Pool
	var v interface{} = &struct{}{}
	pool.Put(v)
}

// Invalid: Put on a pool pointer receiver with non-pointer value.
func badPutOnPointerReceiver() {
	pool := &sync.Pool{}
	pool.Put(123) // MATCH /non-pointer value stored in sync.Pool causes allocation; use a pointer instead/
}

// Valid: Put on a pool pointer receiver with pointer value.
func goodPutOnPointerReceiver() {
	pool := &sync.Pool{}
	x := 123
	pool.Put(&x)
}
