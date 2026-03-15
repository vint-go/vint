package fixtures

// Invalid: top-level type shadows built-in int type
type int struct{} // MATCH /top-level declaration shadows predeclared identifier int/

// Invalid: top-level var shadows built-in error type
var error = "something" // MATCH /top-level declaration shadows predeclared identifier error/

// Invalid: top-level const shadows built-in true constant
const true = 0 // MATCH /top-level declaration shadows predeclared identifier true/

// Invalid: top-level function shadows built-in len function
func len() {} // MATCH /top-level declaration shadows predeclared identifier len/

// Invalid: top-level var shadows built-in nil
var nil = 0 // MATCH /top-level declaration shadows predeclared identifier nil/

// Invalid: top-level type shadows built-in string type
type string = []byte // MATCH /top-level declaration shadows predeclared identifier string/

// Valid: custom type name that does not shadow
type myInt struct{}

// Valid: custom var name that does not shadow
var errMsg = "something"

// Valid: method named with a builtin name is okay
type Foo struct{}

func (f Foo) len() int { return 0 }

// Valid: regular function with non-builtin name
func doSomething() {}
