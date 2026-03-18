package fixtures

// Invalid: all parameters share the same type and could be combined
func foo(a, b int, c, d int, e, f int, g int) {} // MATCH /parameters could be combined by type, e.g. a, b, c, d, e, f, g int/

// Invalid: simple case with two params of same type
func baz(a string, b string) {} // MATCH /parameters could be combined by type, e.g. a, b string/

// Valid: parameters already combined
func good(a, b, c, d, e, f, g int) {}

// Valid: different types for each param
func different(a int, b string, c float64) {}

// Valid: single parameter
func single(a int) {}

// Valid: no parameters
func none() {}

// Valid: multi-line parameter declaration (should be skipped)
func multiLine(
	a int,
	b int,
) {
}

// Valid: unnamed parameters (should be skipped)
type myInterface interface {
	Method(int, int)
}
