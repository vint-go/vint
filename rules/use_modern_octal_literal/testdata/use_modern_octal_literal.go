package fixtures

// Invalid: old-style octal literals.

func oldStyleOctal() {
	mode := 0755  // MATCH /old-style octal literal 0755, use 0o755 instead/
	perm := 0644  // MATCH /old-style octal literal 0644, use 0o644 instead/
	val := 0123   // MATCH /old-style octal literal 0123, use 0o123 instead/
	small := 07   // MATCH /old-style octal literal 07, use 0o7 instead/
	_ = mode
	_ = perm
	_ = val
	_ = small
}

// Valid: modern octal literals.

func modernOctal() {
	mode := 0o755
	perm := 0o644
	val := 0o123
	_ = mode
	_ = perm
	_ = val
}

// Valid: non-octal integer literals.

func nonOctal() {
	x := 42
	y := 0xFF
	z := 0b1010
	w := 0
	_ = x
	_ = y
	_ = z
	_ = w
}
