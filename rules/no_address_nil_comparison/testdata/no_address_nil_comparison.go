package fixtures

func invalidEqualNil() {
	var x int
	if &x == nil { // MATCH /address of a variable is never nil/
		// unreachable
	}
}

func invalidNotEqualNil() {
	var x int
	if &x != nil { // MATCH /address of a variable is never nil/
		// always true
	}
}

func invalidNilOnLeft() {
	var x int
	if nil == &x { // MATCH /address of a variable is never nil/
		// unreachable
	}
}

func invalidNilNotEqualOnLeft() {
	var x int
	if nil != &x { // MATCH /address of a variable is never nil/
		// always true
	}
}

func invalidStructField() {
	type S struct{ field int }
	var s S
	if &s.field == nil { // MATCH /address of a variable is never nil/
		// unreachable
	}
}

func invalidArrayElement() {
	var arr [3]int
	if &arr[0] == nil { // MATCH /address of a variable is never nil/
		// unreachable
	}
}

// Valid: checking a pointer parameter against nil is meaningful
func validPointerParam(x *int) {
	if x == nil {
		return
	}
}

// Valid: comparing two pointers
func validPointerComparison(x, y *int) bool {
	return x == y
}

// Valid: comparing a pointer variable to nil
func validPointerVar() {
	var p *int
	if p == nil {
		return
	}
}

// Valid: address-of in non-comparison context
func validAddressOfUsage() {
	var x int
	p := &x
	_ = p
}
