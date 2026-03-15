package fixtures

// Invalid: struct is 24 bytes due to padding, could be 16 bytes.

type Inefficient struct { // MATCH /struct Inefficient could have its fields rearranged to use 16 bytes instead of 24 bytes (optimal field order: b, a, c)/
	a bool
	b int64
	c bool
}

// Invalid: another struct with suboptimal field order.

type BadOrder struct { // MATCH /struct BadOrder could have its fields rearranged to use 16 bytes instead of 24 bytes (optimal field order: y, x, z)/
	x bool
	y int64
	z bool
}

// Valid: struct is already optimally ordered.

type Efficient struct {
	b int64
	a bool
	c bool
}

// Valid: single field struct has nothing to reorder.

type SingleField struct {
	x int64
}

// Valid: struct with fields of the same size has no padding waste.

type SameSize struct {
	a int64
	b int64
	c int64
}

// Valid: empty struct.

type Empty struct{}
