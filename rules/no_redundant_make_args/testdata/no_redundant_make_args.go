package fixtures

func redundantSliceCapacity() {
	s := make([]int, 10, 10) // MATCH /redundant capacity argument in make, length and capacity are the same/
	_ = s
}

func redundantSliceCapacityWithVariable() {
	n := 5
	s := make([]int, n, n) // MATCH /redundant capacity argument in make, length and capacity are the same/
	_ = s
}

func redundantMapSizeHint() {
	m := make(map[string]int, 0) // MATCH /redundant size hint 0 in make, can be omitted/
	_ = m
}

func validSliceMake() {
	s := make([]int, 10)
	_ = s
}

func validSliceMakeDifferentCap() {
	s := make([]int, 0, 10)
	_ = s
}

func validMapMake() {
	m := make(map[string]int)
	_ = m
}

func validMapMakeWithSize() {
	m := make(map[string]int, 100)
	_ = m
}

func validChannelMake() {
	ch := make(chan int)
	_ = ch
}

func validChannelMakeWithBuffer() {
	ch := make(chan int, 10)
	_ = ch
}
