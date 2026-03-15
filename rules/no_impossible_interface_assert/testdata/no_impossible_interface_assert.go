package fixtures

type Reader interface {
	Read(p []byte) (int, error)
}

type BadReader interface {
	Read(p []byte) error // different signature than Reader.Read
}

type ReadCloser interface {
	Read(p []byte) (int, error)
	Close() error
}

type Writer interface {
	Write(p []byte) (int, error)
}

type BadWriter interface {
	Write(data string) error // different signature than Writer.Write
}

// Invalid: impossible type assertions

func impossibleDirectAssertion(r Reader) {
	_ = r.(BadReader) // MATCH /impossible type assertion: no type can implement both interfaces (method Read has conflicting signatures)/
}

func impossibleCommaOk(r Reader) {
	_, ok := r.(BadReader) // MATCH /impossible type assertion: no type can implement both interfaces (method Read has conflicting signatures)/
	_ = ok
}

func impossibleTypeSwitch(r Reader) {
	switch r.(type) {
	case BadReader: // MATCH /impossible type assertion: no type can implement both interfaces (method Read has conflicting signatures)/
	}
}

func impossibleTypeSwitchAssign(r Reader) {
	switch v := r.(type) {
	case BadReader: // MATCH /impossible type assertion: no type can implement both interfaces (method Read has conflicting signatures)/
		_ = v
	}
}

func impossibleWriterAssertion(w Writer) {
	_ = w.(BadWriter) // MATCH /impossible type assertion: no type can implement both interfaces (method Write has conflicting signatures)/
}

// Valid: possible type assertions

func validReadCloser(r Reader) {
	if rc, ok := r.(ReadCloser); ok {
		defer rc.Close()
	}
}

func validEmptyInterface(x interface{}) {
	if r, ok := x.(Reader); ok {
		_ = r
	}
}

func validNoConflict(r Reader) {
	if w, ok := r.(Writer); ok {
		_ = w
	}
}

func validTypeSwitch(r Reader) {
	switch r.(type) {
	case ReadCloser:
	}
}
