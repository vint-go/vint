package fixtures

// Invalid: stub doc comments on exported declarations.

// ProcessData ...
func ProcessData() {} // MATCH /doc comment for ProcessData appears to be a stub/

// MyType whatever
type MyType struct{} // MATCH /doc comment for MyType appears to be a stub/

// ExportedConst xxx
const ExportedConst = 42 // MATCH /doc comment for ExportedConst appears to be a stub/

// ExportedVar .
var ExportedVar int // MATCH /doc comment for ExportedVar appears to be a stub/

// AnotherFunc todo
func AnotherFunc() {} // MATCH /doc comment for AnotherFunc appears to be a stub/

// StubType -
type StubType int // MATCH /doc comment for StubType appears to be a stub/

// Valid: proper doc comments.

// ProcessGoodData reads input data and transforms it into the output format.
func ProcessGoodData() {}

// GoodType represents a domain entity with associated metadata.
type GoodType struct{}

// GoodConst is the default value used in configuration.
const GoodConst = 100

// GoodVar holds the global state for the application.
var GoodVar string

// Valid: unexported symbols are not checked.

// unexported ...
func unexported() {}

// unexportedType whatever
type unexportedType struct{}
