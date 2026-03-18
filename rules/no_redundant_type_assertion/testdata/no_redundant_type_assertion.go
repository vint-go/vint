package fixtures

import "io"

func getReader() io.Reader {
	return nil
}

func redundantTypeAssertions() {
	// Invalid: asserting to the same interface type
	var r io.Reader = getReader()
	_ = r.(io.Reader) // MATCH /redundant type assertion: expression already has type io.Reader/

	// Valid: asserting to a different interface type
	var r2 io.Reader = getReader()
	_ = r2.(io.ReadCloser)

	// Valid: no assertion, just assignment
	var r3 io.Reader = getReader()
	_ = r3
}
