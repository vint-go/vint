package fixtures

func process(v int) {
	_ = v
}

// Invalid: ranging over a large array copies the entire array.
func badRangeLargeArray() {
	var data [1024]int
	for _, v := range data { // MATCH /range expression copies 8192 bytes, consider using a pointer (e.g. &expr)/
		process(v)
	}
}

// Invalid: ranging over a large byte array.
func badRangeLargeByteArray() {
	var buf [512]byte
	for _, b := range buf { // MATCH /range expression copies 512 bytes, consider using a pointer (e.g. &expr)/
		_ = b
	}
}

// Valid: using a pointer to the array avoids the copy.
func goodRangePointer() {
	var data [1024]int
	for _, v := range &data {
		process(v)
	}
}

// Valid: slices do not have this issue.
func goodRangeSlice() {
	var data []int
	for _, v := range data {
		process(v)
	}
}

// Valid: small array below the default threshold of 512 bytes.
func goodRangeSmallArray() {
	var data [10]int
	for _, v := range data {
		process(v)
	}
}

// Valid: ranging over a map.
func goodRangeMap() {
	m := map[string]int{"a": 1}
	for k, v := range m {
		_ = k
		_ = v
	}
}

// Valid: ranging over a string.
func goodRangeString() {
	s := "hello"
	for _, c := range s {
		_ = c
	}
}

// Valid: ranging over a channel.
func goodRangeChannel() {
	ch := make(chan int, 10)
	go func() { close(ch) }()
	for v := range ch {
		_ = v
	}
}
