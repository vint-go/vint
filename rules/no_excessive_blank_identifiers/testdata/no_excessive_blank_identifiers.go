package fixtures

func myFunc() (int, int, int, int, int) {
	return 1, 2, 3, 4, 5
}

func fourReturnValues() (int, int, int, int) {
	return 1, 2, 3, 4
}

func complexFunction() (int, int, int, int, int, int) {
	return 1, 2, 3, 4, 5, 6
}

func singleReturn() int {
	return 1
}

func invalid() {
	// With default max-blank-identifiers: 2, this has 3 blank identifiers
	a, _, _, _, b := myFunc() // MATCH /assignment has too many blank identifiers (3 > 2)/
	_ = a
	_ = b

	// All return values except one are discarded
	x, _, _, _ := fourReturnValues() // MATCH /assignment has too many blank identifiers (3 > 2)/
	_ = x

	// Excessive blank identifiers in a short variable declaration
	result, _, _, _, _, _ := complexFunction() // MATCH /assignment has too many blank identifiers (5 > 2)/
	_ = result
}

func valid() {
	// Only 1 blank identifier (within default threshold of 2)
	a, _ := myFunc()
	_ = a

	// Exactly 2 blank identifiers (within default threshold of 2)
	b, _, _ := myFunc()
	_ = b

	// No blank identifiers at all
	c, d, e := myFunc()
	_ = c
	_ = d
	_ = e

	// Single return value assignment
	result := singleReturn()
	_ = result
}
