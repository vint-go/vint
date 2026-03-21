package example

func greet(name string) string { //nolint:unparam
	return "hello"
}

func use() {
	_ = greet("world")
}
