package fixtures

import (
	"bytes"
	"io"
)

func useStringWriterInvalid(w io.Writer) {
	w.Write([]byte("hello world")) // MATCH /use WriteString or io.WriteString instead of Write([]byte("..."))/
}

func useStringWriterInvalidBuffer() {
	var buf bytes.Buffer
	buf.Write([]byte("test data")) // MATCH /use WriteString or io.WriteString instead of Write([]byte("..."))/
}

func useStringWriterValid(w io.Writer, name string) {
	// Using io.WriteString directly is fine
	io.WriteString(w, "hello world")

	// Writing a variable byte slice is fine
	data := []byte("test")
	w.Write(data)

	// Writing a non-literal []byte conversion is fine
	w.Write([]byte(name))

	// Writing without []byte conversion is fine
	var buf bytes.Buffer
	buf.WriteString("hello")
}
