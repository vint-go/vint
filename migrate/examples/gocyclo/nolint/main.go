package example

func process(x int) int { //nolint:gocyclo
	if x > 0 {
		return x
	}
	return -x
}
