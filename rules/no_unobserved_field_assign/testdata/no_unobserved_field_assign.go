package fixtures

type Counter struct {
	count int
}

// Invalid: value receiver, changes to count are lost.
func (c Counter) Increment() {
	c.count++ // MATCH /field assignment to value receiver will not be observed; did you mean to use a pointer receiver?/
}

// Invalid: value receiver, field assignment is lost.
func (c Counter) SetCount(n int) {
	c.count = n // MATCH /field assignment to value receiver will not be observed; did you mean to use a pointer receiver?/
}

// Valid: pointer receiver, changes are preserved.
func (c *Counter) IncrementPtr() {
	c.count++
}

// Valid: pointer receiver, assignment is preserved.
func (c *Counter) SetCountPtr(n int) {
	c.count = n
}

type Pair struct {
	x int
	y int
}

// Valid: method returns the receiver, so modification is intentional.
func (p Pair) WithX(x int) Pair {
	p.x = x
	return p
}

// Valid: method returns a field of the receiver.
func (p Pair) GetX() int {
	return p.x
}

// Valid: no field assignment on the receiver.
func (c Counter) Value() int {
	return c.count
}

// Valid: no field assignment at all.
func (c Counter) String() string {
	return "counter"
}
