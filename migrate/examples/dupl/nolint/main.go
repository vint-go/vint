package example

func processA(a int) int {
	x := a + 1
	if x > 0 {
		return x
	}
	return 0
}

func processB(b int) int { //nolint:dupl
	y := b + 1
	if y > 0 {
		return y
	}
	return 0
}
