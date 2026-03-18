package fixtures

type Foo struct {
	name string
}

func (f Foo) Bar() string {
	return f.name
}

func (f *Foo) SetName(name string) {
	f.name = name
}

func (f Foo) Add(a, b int) int {
	return a + b
}

func example() {
	f := Foo{name: "test"}

	// Invalid: method expression calls
	Foo.Bar(f)                // MATCH /method expression call can be replaced with f.Bar()/
	(*Foo).SetName(&f, "new") // MATCH /method expression call can be replaced with &f.SetName()/
	Foo.Add(f, 1, 2)         // MATCH /method expression call can be replaced with f.Add()/

	// Valid: direct method calls
	f.Bar()
	f.SetName("new")
	f.Add(1, 2)
}
