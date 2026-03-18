package fixtures

// Represents a server configuration.
type Config struct { // MATCH /comment on exported type Config should be of the form "Config ..."/
	Host string
	Port int
}

// Config2 represents a server configuration.
type Config2 struct {
	Host string
	Port int
}

// unexported types should not be flagged
type unexportedType struct{}

// Handler is an HTTP handler.
type Handler struct{}

// Does something important.
type Worker struct{} // MATCH /comment on exported type Worker should be of the form "Worker ..."/

type NoDocType struct{}

// Runner starts background tasks.
type Runner interface{}

// Starts background processes.
type Daemon interface{} // MATCH /comment on exported type Daemon should be of the form "Daemon ..."/

// Validator checks input data.
type Validator func(string) bool

// Checks input data.
type Checker func(string) bool // MATCH /comment on exported type Checker should be of the form "Checker ..."/

// Counter is an integer counter type.
type Counter int

// An integer type.
type Number int // MATCH /comment on exported type Number should be of the form "Number ..."/
