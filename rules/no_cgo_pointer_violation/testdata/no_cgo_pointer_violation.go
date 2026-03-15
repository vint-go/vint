package fixtures

/*
#include <stdlib.h>

void process(void *p) {}
void doWork(void *p) {}
*/
import "C"
import "unsafe"

func invalidStringPointer() {
	s := "hello"
	// Passing a pointer to a Go string (which contains a Go pointer)
	C.process(unsafe.Pointer(&s)) // MATCH /possible violation of cgo pointer passing rules: Go pointer passed to C contains nested Go pointers/
}

func invalidSlicePointer() {
	data := []int{1, 2, 3}
	// Passing a pointer to a slice header (which contains a Go pointer)
	C.process(unsafe.Pointer(&data)) // MATCH /possible violation of cgo pointer passing rules: Go pointer passed to C contains nested Go pointers/
}

func invalidMapPointer() {
	m := map[string]int{"a": 1}
	C.process(unsafe.Pointer(&m)) // MATCH /possible violation of cgo pointer passing rules: Go pointer passed to C contains nested Go pointers/
}

func invalidInterfacePointer() {
	var iface interface{} = 42
	C.process(unsafe.Pointer(&iface)) // MATCH /possible violation of cgo pointer passing rules: Go pointer passed to C contains nested Go pointers/
}

func validByteSliceElement() {
	b := make([]byte, 5)
	copy(b, "hello")
	// Passing a pointer to a byte slice's underlying array element (no nested Go pointers)
	C.process(unsafe.Pointer(&b[0]))
}

func validIntVariable() {
	x := 42
	// Passing a pointer to a plain int (no nested Go pointers)
	C.process(unsafe.Pointer(&x))
}

func noCgoCall() {
	s := "hello"
	_ = unsafe.Pointer(&s) // not passed to C, no violation
}
