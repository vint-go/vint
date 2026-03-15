package fixtures

import "unsafe"

func invalidStoredUintptr() {
	var x int
	// Bad: storing uintptr in a variable and converting later
	// The GC may move the object between the two lines
	ptr := uintptr(unsafe.Pointer(&x))
	p := unsafe.Pointer(ptr) // MATCH /possibly invalid conversion of uintptr to unsafe.Pointer/
	_ = p
}

func validSingleExpression() {
	var x int
	// Good: conversion in a single expression
	p := unsafe.Pointer(uintptr(unsafe.Pointer(&x)) + unsafe.Sizeof(x))
	_ = p
}

func validPrintAddress() {
	var x int
	// Good: converting unsafe.Pointer to uintptr for printing (not converting back)
	_ = uintptr(unsafe.Pointer(&x))
}

func validPointerToPointer() {
	var x int
	// Good: converting pointer to unsafe.Pointer directly (not uintptr)
	p := unsafe.Pointer(&x)
	_ = p
}

func invalidMultipleVars() {
	var a, b int
	u1 := uintptr(unsafe.Pointer(&a))
	u2 := uintptr(unsafe.Pointer(&b))
	p1 := unsafe.Pointer(u1) // MATCH /possibly invalid conversion of uintptr to unsafe.Pointer/
	p2 := unsafe.Pointer(u2) // MATCH /possibly invalid conversion of uintptr to unsafe.Pointer/
	_ = p1
	_ = p2
}
