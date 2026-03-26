package fixtures

// Invalid: unexported struct type never referenced anywhere.

type config struct { // MATCH /type config is unused/
	host string
	port int
}

func Connect() {}

// Invalid: unexported function type never referenced anywhere.

type handler func(req string) string // MATCH /type handler is unused/

func Process(input string) string {
	return input
}

// Valid: unexported named type referenced by unexported const declarations.

type color int

const (
	red   color = iota
	green
	blue
)

// Valid: unexported struct used in a function return type.

type dbConfig struct {
	host string
	port int
}

func NewConfig() dbConfig {
	return dbConfig{host: "localhost", port: 8080}
}

// Valid: exported types are always considered used.

type Server struct {
	Addr string
}

// Valid: unexported interface used in a function parameter.

type stringer interface {
	String() string
}

func Format(s stringer) string {
	return s.String()
}

// Valid: unexported type used as element type in an unexported variable declaration.

type tableTest struct {
	name     string
	input    string
	expected string
}

var myTests = []tableTest{
	{name: "test1", input: "a", expected: "b"},
}

// Valid: two types referencing each other, both used from a function.

type middleware func(next reqHandler) reqHandler

type reqHandler func(req string) string

func Chain(h reqHandler, mws ...middleware) reqHandler {
	for _, mw := range mws {
		h = mw(h)
	}
	return h
}
