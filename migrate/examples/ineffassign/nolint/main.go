package example

func compute() int {
	x := 1    //nolint:ineffassign
	x = 2
	return x
}
