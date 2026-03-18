package fixtures

import (
	"bytes"
	"fmt"
	"io"
)

func useFprintInvalid(w io.Writer, name string) {
	w.Write([]byte(fmt.Sprintf("hello %s", name))) // MATCH /fmt.Sprintf result converted to []byte and written; use fmt.Fprintf instead/
	w.Write([]byte(fmt.Sprint("hello ", name)))     // MATCH /fmt.Sprint result converted to []byte and written; use fmt.Fprint instead/
}

func useFprintInvalidBuffer(name string) {
	var buf bytes.Buffer
	buf.Write([]byte(fmt.Sprintf("value: %d", 42))) // MATCH /fmt.Sprintf result converted to []byte and written; use fmt.Fprintf instead/
	buf.Write([]byte(fmt.Sprint("test")))            // MATCH /fmt.Sprint result converted to []byte and written; use fmt.Fprint instead/
}

func useFprintValid(w io.Writer, name string) {
	// Using fmt.Fprintf directly is fine
	fmt.Fprintf(w, "hello %s", name)

	// Using fmt.Fprint directly is fine
	fmt.Fprint(w, "hello ", name)

	// Writing a regular byte slice is fine
	w.Write([]byte("hello"))

	// Writing a non-fmt function result is fine
	w.Write([]byte(name))

	// Using Write with a variable is fine
	data := []byte("test")
	w.Write(data)
}
