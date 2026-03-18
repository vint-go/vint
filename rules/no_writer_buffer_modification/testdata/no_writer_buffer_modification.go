package fixtures

// Invalid: modifying the input slice in a Write method.

type badWriter struct{}

func (w badWriter) Write(p []byte) (int, error) {
	for i := range p {
		p[i] = p[i] ^ 0xFF // MATCH /an io.Writer must not modify the provided buffer/
	}
	return len(p), nil
}

// Invalid: assigning to an index of the buffer directly.

type badWriter2 struct{}

func (w badWriter2) Write(p []byte) (int, error) {
	p[0] = 0 // MATCH /an io.Writer must not modify the provided buffer/
	return len(p), nil
}

// Valid: creating a copy and modifying that.

type goodWriter struct{}

func (w goodWriter) Write(p []byte) (int, error) {
	buf := make([]byte, len(p))
	copy(buf, p)
	for i := range buf {
		buf[i] = buf[i] ^ 0xFF
	}
	return len(p), nil
}

// Valid: reading from the buffer without modification.

type goodWriter2 struct{}

func (w goodWriter2) Write(p []byte) (int, error) {
	total := 0
	for _, b := range p {
		total += int(b)
	}
	_ = total
	return len(p), nil
}

// Valid: not a method (no receiver), so not an io.Writer.

func Write(p []byte) (int, error) {
	p[0] = 0
	return len(p), nil
}

// Valid: method named Write but wrong signature (not io.Writer).

type notWriter struct{}

func (w notWriter) Write(data string) error {
	return nil
}

// Valid: pointer receiver Write method with no modification.

type ptrWriter struct{}

func (w *ptrWriter) Write(p []byte) (int, error) {
	return len(p), nil
}
