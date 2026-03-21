package example

func convert() int {
	var x int = 42
	return int(x) //nolint:unconvert
}
