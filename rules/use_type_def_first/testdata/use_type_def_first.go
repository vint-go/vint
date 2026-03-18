package fixtures

// Invalid: method declared before type definition

func (s MyStruct) Process() error { // MATCH /method Process declared before type definition of MyStruct/
	return nil
}

type MyStruct struct {
	Name string
}

// Invalid: pointer receiver before type definition

func (s *AnotherStruct) Run() { // MATCH /method Run declared before type definition of AnotherStruct/
}

type AnotherStruct struct {
	Value int
}

// Valid: type defined first, then methods

type GoodStruct struct {
	Field string
}

func (g GoodStruct) DoSomething() {}

func (g *GoodStruct) DoAnotherThing() {}

// Valid: function without receiver (not a method)

func freeFunction() {}

// Valid: method on type defined in another file (not flagged)

func (e ExternalType) Method() {}
