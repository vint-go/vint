package fixtures

func immediateNewDerefInvalid() {
	x := *new(bool) // MATCH /immediate dereference of new(bool), use zero value instead/
	_ = x

	y := *new(int) // MATCH /immediate dereference of new(int), use zero value instead/
	_ = y

	s := *new(string) // MATCH /immediate dereference of new(string), use zero value instead/
	_ = s

	f := *new(float64) // MATCH /immediate dereference of new(float64), use zero value instead/
	_ = f
}

func immediateNewDerefValid() {
	// Using zero values directly is fine
	x := false
	_ = x

	y := 0
	_ = y

	// Using new without immediate dereference is fine
	p := new(int)
	_ = p
}

// Allowed for generic type parameters
func zero[T any]() T {
	return *new(T)
}
