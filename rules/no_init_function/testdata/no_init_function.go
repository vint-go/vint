package fixtures

// Invalid: init functions should not be used.

func init() { // MATCH /init function found, consider using explicit initialization/
	setupDatabase()
}

func init() { // MATCH /init function found, consider using explicit initialization/
	registerHandlers()
}

// Valid: regular functions are fine.

func Setup() {
	setupDatabase()
}

func Initialize() {
	registerHandlers()
}

// Valid: methods named init on a receiver are fine.

type MyStruct struct{}

func (m *MyStruct) init() {
	// This is a method, not an init function.
}

// helper stubs to make the file parse
func setupDatabase()    {}
func registerHandlers() {}
