package fixtures

type TypeA struct{}
type TypeB struct{}
type TypeC struct{}
type TypeD struct{}

func handleA(v *TypeA) {}
func handleB(v *TypeB) {}
func handleC(v *TypeC) {}
func handleD(v *TypeD) {}
func handleDefault()   {}

// Invalid: if-else chain with 3 type assertions on the same variable
func invalidThreeTypeAssertions(x interface{}) {
	if v, ok := x.(*TypeA); ok { // MATCH /could replace type assertion if-else chain on x with type switch/
		handleA(v)
	} else if v, ok := x.(*TypeB); ok {
		handleB(v)
	} else if v, ok := x.(*TypeC); ok {
		handleC(v)
	}
}

// Invalid: if-else chain with 2 type assertions on the same variable
func invalidTwoTypeAssertions(x interface{}) {
	if v, ok := x.(*TypeA); ok { // MATCH /could replace type assertion if-else chain on x with type switch/
		handleA(v)
	} else if v, ok := x.(*TypeB); ok {
		handleB(v)
	}
}

// Invalid: if-else chain with type assertions and a final else
func invalidWithElse(x interface{}) {
	if v, ok := x.(*TypeA); ok { // MATCH /could replace type assertion if-else chain on x with type switch/
		handleA(v)
	} else if v, ok := x.(*TypeB); ok {
		handleB(v)
	} else if v, ok := x.(*TypeC); ok {
		handleC(v)
	} else {
		handleDefault()
	}
}

// Invalid: 4 type assertions
func invalidFourTypeAssertions(x interface{}) {
	if v, ok := x.(*TypeA); ok { // MATCH /could replace type assertion if-else chain on x with type switch/
		handleA(v)
	} else if v, ok := x.(*TypeB); ok {
		handleB(v)
	} else if v, ok := x.(*TypeC); ok {
		handleC(v)
	} else if v, ok := x.(*TypeD); ok {
		handleD(v)
	}
}

// Valid: already uses a type switch
func validTypeSwitch(x interface{}) {
	switch v := x.(type) {
	case *TypeA:
		handleA(v)
	case *TypeB:
		handleB(v)
	case *TypeC:
		handleC(v)
	}
}

// Valid: single type assertion (not a chain)
func validSingleTypeAssertion(x interface{}) {
	if v, ok := x.(*TypeA); ok {
		handleA(v)
	}
}

// Valid: type assertion with else but not else-if
func validTypeAssertionWithElse(x interface{}) {
	if v, ok := x.(*TypeA); ok {
		handleA(v)
	} else {
		handleDefault()
	}
}

// Valid: if-else chain but not type assertions
func validNonTypeAssertionChain(x int) {
	if x == 1 {
		handleDefault()
	} else if x == 2 {
		handleDefault()
	}
}

// Valid: type assertions on different variables
func validDifferentVariables(x interface{}, y interface{}) {
	if v, ok := x.(*TypeA); ok {
		handleA(v)
	} else if v, ok := y.(*TypeB); ok {
		handleB(v)
	}
}
