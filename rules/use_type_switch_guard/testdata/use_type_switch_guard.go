package fixtures

type TypeA struct{ Name string }
type TypeB struct{ Value int }
type TypeC struct{ Data string }

func handleA(a *TypeA) {}
func handleB(b *TypeB) {}
func handleC(c *TypeC) {}

// Invalid: type switch without guard with redundant type assertions
func invalidBasic(v interface{}) {
	switch v.(type) { // MATCH /use type switch guard (switch v := v.(type)) to avoid redundant type assertions/
	case *TypeA:
		handleA(v.(*TypeA))
	case *TypeB:
		handleB(v.(*TypeB))
	}
}

// Invalid: type switch with assertions in multiple case clauses
func invalidMultipleCases(v interface{}) {
	switch v.(type) { // MATCH /use type switch guard (switch v := v.(type)) to avoid redundant type assertions/
	case *TypeA:
		handleA(v.(*TypeA))
	case *TypeB:
		handleB(v.(*TypeB))
	case *TypeC:
		handleC(v.(*TypeC))
	}
}

// Invalid: only one case has a type assertion but that is enough
func invalidSingleCaseAssertion(v interface{}) {
	switch v.(type) { // MATCH /use type switch guard (switch v := v.(type)) to avoid redundant type assertions/
	case *TypeA:
		handleA(v.(*TypeA))
	case *TypeB:
		// no assertion here
	}
}

// Valid: type switch already uses a guard variable
func validWithGuard(v interface{}) {
	switch v := v.(type) {
	case *TypeA:
		handleA(v)
	case *TypeB:
		handleB(v)
	}
}

// Valid: type switch without any type assertions in cases
func validNoAssertions(v interface{}) {
	switch v.(type) {
	case *TypeA:
		// just do something
	case *TypeB:
		// just do something
	}
}

// Valid: multi-type case clause (results in interface{}, skip)
func validMultiTypeCase(v interface{}) {
	switch v.(type) {
	case *TypeA, *TypeB:
		// multi-type case
	}
}
