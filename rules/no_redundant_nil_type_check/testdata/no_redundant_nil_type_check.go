package fixtures

// Invalid: redundant nil check before type assertion.
func bad1(x interface{}) {
	if x != nil { // MATCH /redundant nil check on variable before type assertion; type assertion already handles nil/
		if v, ok := x.(string); ok {
			_ = v
		}
	}
}

// Invalid: nil != x form (reversed operands).
func bad2(x interface{}) {
	if nil != x { // MATCH /redundant nil check on variable before type assertion; type assertion already handles nil/
		if v, ok := x.(int); ok {
			_ = v
		}
	}
}

// Invalid: type assertion to a struct type.
func bad3(x interface{}) {
	if x != nil { // MATCH /redundant nil check on variable before type assertion; type assertion already handles nil/
		if v, ok := x.(myStruct); ok {
			_ = v
		}
	}
}

// Valid: the inner if body does more than just a type assertion.
func good1(x interface{}) {
	if x != nil {
		doSomething()
		if v, ok := x.(string); ok {
			_ = v
		}
	}
}

// Valid: the type assertion is on a different variable.
func good2(x interface{}, y interface{}) {
	if x != nil {
		if v, ok := y.(string); ok {
			_ = v
		}
	}
}

// Valid: direct type assertion without ok check (single-value).
func good3(x interface{}) {
	if x != nil {
		v := x.(string)
		_ = v
	}
}

// Valid: no type assertion at all.
func good4(x interface{}) {
	if x != nil {
		doSomething()
	}
}

// Valid: outer if has an else branch.
func good5(x interface{}) {
	if x != nil {
		if v, ok := x.(string); ok {
			_ = v
		}
	} else {
		doSomething()
	}
}

// Valid: condition is not a nil check.
func good6(x interface{}) {
	if x == nil {
		if v, ok := x.(string); ok {
			_ = v
		}
	}
}

// Valid: outer if has init statement.
func good7(x interface{}) {
	if y := x; y != nil {
		if v, ok := y.(string); ok {
			_ = v
		}
	}
}

// Valid: type assertion already used directly with if (no nil check wrapping).
func good8(x interface{}) {
	if v, ok := x.(string); ok {
		_ = v
	}
}

type myStruct struct{}

func doSomething() {}
