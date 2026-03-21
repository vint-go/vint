package example

// go:generate stringer -type=Pill //nolint:gocheckcompilerdirectives

//go:embod hello.txt //nolint:gocheckcompilerdirectives

func Add(a, b int) int {
	return a + b
}
