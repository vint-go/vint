package fixtures

type Point struct {
	X, Y int
}

// Invalid: writing to copy that is never used
func unusedWriteToCopy() Point {
	p := Point{X: 1, Y: 2}
	q := p
	q.X = 10 // MATCH /unused write to field or element of q/
	return p
}

// Invalid: writing to struct field where struct is never used after
func unusedWriteToLocal() {
	p := Point{X: 1, Y: 2}
	p.X = 10 // MATCH /unused write to field or element of p/
}

// Invalid: writing to array element where array is never used after
func unusedWriteToArray() {
	var arr [3]int
	arr[0] = 42 // MATCH /unused write to field or element of arr/
}

// Valid: writing to p which is then returned
func validWriteAndReturn() Point {
	p := Point{X: 1, Y: 2}
	p.X = 10
	return p
}

// Valid: modified copy is used (returned)
func validCopyUsed() Point {
	p := Point{X: 1, Y: 2}
	q := p
	q.X = 10
	return q
}

// Valid: variable is passed to a function after the write
func validPassedToFunc() {
	p := Point{X: 1, Y: 2}
	p.X = 10
	usePoint(p)
}

// Valid: variable used in a later expression
func validUsedInExpression() int {
	p := Point{X: 1, Y: 2}
	p.X = 10
	return p.X + p.Y
}

// Valid: writing to a slice element (reference type, not value type)
func validSliceWrite() {
	s := []int{1, 2, 3}
	s[0] = 42
}

func usePoint(_ Point) {}
