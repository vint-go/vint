package fixtures

// Invalid: wrong signature for fmt.Stringer (should return string, not error)

type BadStringer struct{}

func (m BadStringer) String() error { // MATCH /method String has signature mismatch with fmt.Stringer interface/
	return nil
}

// Invalid: wrong signature for encoding/json.Marshaler (should return ([]byte, error))

type BadMarshaler struct{}

func (m BadMarshaler) MarshalJSON() error { // MATCH /method MarshalJSON has signature mismatch with encoding/json.Marshaler interface/
	return nil
}

// Invalid: wrong signature for io.Reader (should take []byte and return (int, error))

type BadReader struct{}

func (m BadReader) Read() error { // MATCH /method Read has signature mismatch with io.Reader interface/
	return nil
}

// Invalid: wrong signature for io.Writer (wrong return type)

type BadWriter struct{}

func (m BadWriter) Write(p []byte) error { // MATCH /method Write has signature mismatch with io.Writer interface/
	return nil
}

// Invalid: wrong signature for io.Closer (should return error, not int)

type BadCloser struct{}

func (m BadCloser) Close() int { // MATCH /method Close has signature mismatch with io.Closer interface/
	return 0
}

// Invalid: wrong signature for encoding/json.Unmarshaler (should take []byte)

type BadUnmarshaler struct{}

func (m BadUnmarshaler) UnmarshalJSON() error { // MATCH /method UnmarshalJSON has signature mismatch with encoding/json.Unmarshaler interface/
	return nil
}

// Invalid: wrong signature for error interface (should return string, not int)

type BadError struct{}

func (m BadError) Error() int { // MATCH /method Error has signature mismatch with error interface/
	return 0
}

// Valid: correct signature for fmt.Stringer

type GoodStringer struct{}

func (m GoodStringer) String() string {
	return "GoodStringer"
}

// Valid: correct signature for encoding/json.Marshaler

type GoodMarshaler struct{}

func (m GoodMarshaler) MarshalJSON() ([]byte, error) {
	return []byte(`"good"`), nil
}

// Valid: correct signature for io.Reader

type GoodReader struct{}

func (m GoodReader) Read(p []byte) (int, error) {
	return 0, nil
}

// Valid: correct signature for io.Writer

type GoodWriter struct{}

func (m GoodWriter) Write(p []byte) (int, error) {
	return 0, nil
}

// Valid: correct signature for io.Closer

type GoodCloser struct{}

func (m GoodCloser) Close() error {
	return nil
}

// Valid: correct signature for encoding/json.Unmarshaler

type GoodUnmarshaler struct{}

func (m GoodUnmarshaler) UnmarshalJSON(data []byte) error {
	return nil
}

// Valid: correct signature for error interface

type GoodError struct{}

func (m GoodError) Error() string {
	return "error"
}

// Valid: not a method (plain function), should not be flagged

func String() error {
	return nil
}
