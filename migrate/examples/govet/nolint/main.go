package example

func doSomething() {
	x := 1
	x = x //nolint:govet
	_ = x
}
