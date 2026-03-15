package fixtures

import "fmt"

// Invalid: dereferencing p after confirming it is nil
func nilDerefAfterCheck(p *int) int {
	if p == nil {
		return *p // MATCH /nil pointer dereference: p is nil in this branch/
	}
	return 0
}

// Invalid: accessing a field on nil pointer after confirming nil
func nilFieldAccessAfterCheck(p *struct{ Field int }) int {
	if p == nil {
		return p.Field // MATCH /nil pointer dereference: p is nil in this branch/
	}
	return 0
}

// Invalid: degenerate nil comparison - p is always nil
func degenerateNilComparison() {
	var p *int
	if p != nil { // MATCH /degenerate nil comparison: p is always nil/
		fmt.Println(*p)
	}
}

// Valid: proper nil check followed by return
func validNilCheck(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}

// Valid: pointer is assigned before comparison
func validAssignedPointer() {
	var p *int
	x := 42
	p = &x
	if p != nil {
		fmt.Println(*p)
	}
}

// Valid: non-pointer variable
func validNonPointer() {
	x := 42
	if x != 0 {
		fmt.Println(x)
	}
}
