package fixtures

// Invalid: assembly function Add has wrong argument size (should be 24 for two int64 args + int64 return).
func Add(x, y int64) int64 // MATCH /assembly function Add has argument size 16, but Go declaration expects 24/

// Valid: assembly function Sub has correct argument size.
func Sub(x, y int64) int64

// Valid: regular function with body, not an assembly stub.
func Regular(a, b int) int {
	return a + b
}
