package fixtures

// Invalid: const group where only the first has an explicit value
const (
	A = 1
	B // MATCH /constant implicitly repeats the previous value; consider making it explicit/
	C // MATCH /constant implicitly repeats the previous value; consider making it explicit/
)

// Invalid: string constant with implicit repetition
const (
	S1 = "hello"
	S2 // MATCH /constant implicitly repeats the previous value; consider making it explicit/
)

// Invalid: first has explicit value, second does not, third has explicit
const (
	X1 = 10
	X2     // MATCH /constant implicitly repeats the previous value; consider making it explicit/
	X3 = 30
)

// Valid: all constants have explicit values
const (
	V1 = 1
	V2 = 2
	V3 = 3
)

// Valid: iota-based pattern
const (
	I0 = iota
	I1
	I2
)

// Valid: iota with expression
const (
	Bit0 = 1 << iota
	Bit1
	Bit2
)

// Valid: iota in arithmetic expression
const (
	Off0 = iota + 10
	Off1
	Off2
)

// Valid: single constant (not a group)
const Single = 42

// Valid: single constant in group
const (
	OnlyOne = 100
)

// Valid: no explicit values at all (first spec also has no value — not possible in valid Go, skip)

// Valid: all have values
const (
	All1 = "a"
	All2 = "b"
	All3 = "c"
)
