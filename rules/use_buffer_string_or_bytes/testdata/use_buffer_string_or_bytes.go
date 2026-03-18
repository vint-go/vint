package fixtures

import "bytes"

func invalidStringConversion(buf bytes.Buffer) string {
	return string(buf.Bytes()) // MATCH /use buf.String() instead of string(buf.Bytes())/
}

func invalidStringConversionPointer(buf *bytes.Buffer) string {
	return string(buf.Bytes()) // MATCH /use buf.String() instead of string(buf.Bytes())/
}

func validStringMethod(buf bytes.Buffer) string {
	return buf.String()
}

func validStringMethodPointer(buf *bytes.Buffer) string {
	return buf.String()
}

func validOtherConversion() string {
	b := []byte("hello")
	return string(b)
}
