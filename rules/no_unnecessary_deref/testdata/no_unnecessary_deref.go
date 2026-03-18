package fixtures

type myStruct struct {
	field  int
	nested *myStruct
}

func (s *myStruct) pointerMethod() int {
	return s.field
}

func (s myStruct) valueMethod() int {
	return s.field
}

// Invalid: unnecessary dereference on field access
func badFieldAccess(k *myStruct) {
	_ = (*k).field // MATCH /expression can be simplified to k.field/
}

// Invalid: unnecessary dereference on nested field access
func badNestedFieldAccess(k *myStruct) {
	(*k).field = 42 // MATCH /expression can be simplified to k.field/
}

// Invalid: unnecessary dereference on array index
func badArrayIndex(a *[5]int) {
	_ = (*a)[3] // MATCH /expression can be simplified to a[3]/
}

// Valid: direct field access without deref
func goodFieldAccess(k *myStruct) {
	_ = k.field
}

// Valid: direct array index without deref
func goodArrayIndex(a *[5]int) {
	_ = a[3]
}

// Valid: dereference on pointer receiver method call (skipRecvDeref=true by default)
func goodRecvDeref(k *myStruct) {
	_ = (*k).pointerMethod()
}

// Invalid: dereference on value receiver method call
func badValueMethodDeref(k *myStruct) {
	_ = (*k).valueMethod() // MATCH /expression can be simplified to k.valueMethod/
}

// Valid: dereference on non-pointer type (not a real pointer deref situation)
func goodNonPointer() {
	var x int
	_ = x
}

// Invalid: unnecessary dereference with variable name
func badFieldAccessVar(s *myStruct) {
	_ = (*s).field // MATCH /expression can be simplified to s.field/
}
