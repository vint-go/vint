package fixtures

func useAssignmentOperatorInvalid() {
	x := 0
	x = x + 1 // MATCH /replace x = x + 1 with x += 1/

	y := 1
	y = y * 2 // MATCH /replace y = y * 2 with y *= 2/

	count := 10
	delta := 3
	count = count - delta // MATCH /replace count = count - delta with count -= delta/

	z := 100
	z = z / 5 // MATCH /replace z = z / 5 with z /= 5/

	m := 17
	m = m % 3 // MATCH /replace m = m % 3 with m %= 3/

	bits := 0xFF
	bits = bits & 0x0F // MATCH /replace bits = bits & 0x0F with bits &= 0x0F/

	flags := 0
	flags = flags | 0x01 // MATCH /replace flags = flags | 0x01 with flags |= 0x01/

	val := 0xAB
	val = val ^ 0xFF // MATCH /replace val = val ^ 0xFF with val ^= 0xFF/

	shift := 1
	shift = shift << 2 // MATCH /replace shift = shift << 2 with shift <<= 2/

	big := 1024
	big = big >> 1 // MATCH /replace big = big >> 1 with big >>= 1/

	mask := 0xFF
	mask = mask &^ 0x0F // MATCH /replace mask = mask &^ 0x0F with mask &^= 0x0F/

	_ = x
	_ = y
	_ = count
	_ = delta
	_ = z
	_ = m
	_ = bits
	_ = flags
	_ = val
	_ = shift
	_ = big
	_ = mask
}

func useAssignmentOperatorValid() {
	// Already using assignment operators
	x := 0
	x += 1
	x -= 2
	x *= 3
	x /= 4
	x %= 5
	_ = x

	// Different variable on LHS and RHS
	a := 1
	b := 2
	a = b + 1
	_ = a
	_ = b

	// Short variable declaration
	c := 5 + 3
	_ = c

	// Multiple assignments
	d, e := 1, 2
	d, e = e, d
	_ = d
	_ = e

	// The variable appears on the right side of the binary expression, not the left
	f := 10
	g := 3
	f = g + f
	_ = f
	_ = g
}
