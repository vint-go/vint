package fixtures

type processor struct {
	positives []int
	negatives []int
}

// Invalid: result of append assigned to a different slice
func mismatchedAppendField(p *processor) {
	x := 1
	p.positives = append(p.negatives, x) // MATCH /append result is assigned to a different slice 'p.positives' instead of 'p.negatives'/
}

func mismatchedAppendSimple() {
	a := []int{1, 2}
	b := []int{3, 4}
	item := 5
	a = append(b, item) // MATCH /append result is assigned to a different slice 'a' instead of 'b'/
}

// Valid: result of append assigned to the same slice
func matchedAppendField(p *processor) {
	x := 1
	p.positives = append(p.positives, x)
}

func matchedAppendSimple() {
	a := []int{1, 2}
	a = append(a, 3)
}

// Valid: using ellipsis pattern
func appendWithEllipsis() {
	xs := []int{1, 2}
	ys := []int{3, 4}
	xs = append(ys, xs...)
}

// Valid: blank identifier assignment
func blankAppend() {
	a := []int{1, 2}
	_ = append(a, 3)
}

// Valid: append with multiple elements to the same slice
func appendMultiple() {
	s := []int{1}
	s = append(s, 2, 3, 4)
}
