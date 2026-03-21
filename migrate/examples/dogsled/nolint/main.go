package example

func multiReturn() (int, int, int, int, int) {
	return 1, 2, 3, 4, 5
}

func doSomething() {
	x, _, _, _, _ := multiReturn() //nolint:dogsled
	_ = x
}
