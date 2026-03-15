package fixtures

import "io"

type Counter int64

type ID int

type Constructor func() ID

func unnecessaryConversions() {
	// Invalid: unnecessary conversion of an int variable to int
	var x int = 42
	_ = int(x) // MATCH /unnecessary conversion to int/

	// Invalid: unnecessary conversion of a string variable to string
	var name string = "hello"
	_ = string(name) // MATCH /unnecessary conversion to string/

	// Invalid: unnecessary conversion of a bool variable to bool
	var flag bool = true
	_ = bool(flag) // MATCH /unnecessary conversion to bool/

	// Invalid: unnecessary conversion of a custom type to itself
	var c Counter = 10
	_ = Counter(c) // MATCH /unnecessary conversion to Counter/

	// Invalid: unnecessary pointer conversion
	var p *int
	_ = (*int)(p) // MATCH /unnecessary conversion to (*int)/

	// Invalid: unnecessary conversion in an interface type
	var w io.Writer
	_ = io.Writer(w) // MATCH /unnecessary conversion to io.Writer/

	// Invalid: unnecessary conversion of a function type
	var ctor Constructor
	_ = Constructor(ctor) // MATCH /unnecessary conversion to Constructor/

	// Valid: converting between different types is necessary
	var x32 int32 = 42
	_ = int64(x32)

	// Valid: converting an untyped constant to a specific type
	const maxSize = 100
	_ = int64(maxSize)

	// Valid: float-to-float conversion preserves rounding semantics (without fast-math)
	var f1 float64 = 1.1
	var f2 float64 = 2.2
	_ = float64(f1 * f2)

	// Valid: converting between a named type and its underlying type
	type UserID string
	var uname string = "alice"
	_ = UserID(uname)

	// Valid: converting from one numeric type to another
	var count uint32 = 5
	_ = int(count)

	// Valid: complex number conversion preserves precision semantics (without fast-math)
	var c1 complex128 = complex(1.0, 2.0)
	_ = complex128(c1)
}
