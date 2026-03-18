package fixtures

type T struct{ Name string }

func addressOfDerefInvalid(t *T) *T {
	return &*t // MATCH /&*x does not copy x, use a local variable to copy the value/
}

func addressOfDerefInvalidAssign(t *T) {
	p := &*t // MATCH /&*x does not copy x, use a local variable to copy the value/
	_ = p
}

func addressOfDerefValid(t *T) *T {
	// Creates a real copy
	v := *t
	return &v
}

func addressOfDerefValidNoDeref(t *T) *T {
	// Simply returning the pointer is fine
	return t
}
