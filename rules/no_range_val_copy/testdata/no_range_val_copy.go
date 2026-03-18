package fixtures

// BigStruct has 17 int64 fields = 136 bytes on amd64, exceeding 128 threshold.
type BigStruct struct {
	A, B, C, D, E int64
	F, G, H, I, J int64
	K, L, M, N, O int64
	P, Q          int64
}

// SmallStruct has 2 int64 fields = 16 bytes, under the threshold.
type SmallStruct struct {
	X, Y int64
}

var bigItems []BigStruct
var smallItems []SmallStruct

// Invalid: range value copies BigStruct (136 bytes) on each iteration.
func badRangeValCopy() {
	for _, item := range bigItems { // MATCH /range value 'item' copies 136 bytes each iteration, consider using index access or taking the address/
		_ = item
	}
}

// Valid: using index access avoids the copy.
func goodIndexAccess() {
	for i := range bigItems {
		_ = bigItems[i]
	}
}

// Valid: only key variable, no value copy.
func goodKeyOnly() {
	for i := range bigItems {
		_ = i
	}
}

// Valid: blank identifier for value, no copy happens.
func goodBlankValue() {
	for _, _ = range bigItems {
	}
}

// Valid: small struct under threshold.
func goodSmallStruct() {
	for _, item := range smallItems {
		_ = item
	}
}

// Valid: pointer slice, value is a pointer (no large copy).
func goodPointerSlice() {
	var items []*BigStruct
	for _, item := range items {
		_ = item
	}
}

// Valid: ranging over a map with small values.
func goodMapSmallValue() {
	m := map[string]int{}
	for _, v := range m {
		_ = v
	}
}

// Valid: ranging over a string (rune values are small).
func goodStringRange() {
	for _, r := range "hello" {
		_ = r
	}
}

// Valid: ranging over a channel of small values.
func goodChannelRange() {
	ch := make(chan int)
	for v := range ch {
		_ = v
	}
}

// Invalid: array of BigStruct, value copy is still 136 bytes.
func badArrayRange() {
	var arr [10]BigStruct
	for _, item := range arr { // MATCH /range value 'item' copies 136 bytes each iteration, consider using index access or taking the address/
		_ = item
	}
}

// BigArray wraps a large fixed-size array: 256 bytes.
type BigArray struct {
	Data [256]byte
}

var bigArrayItems []BigArray

// Invalid: BigArray is 256 bytes, exceeding threshold.
func badBigArrayRange() {
	for _, item := range bigArrayItems { // MATCH /range value 'item' copies 256 bytes each iteration, consider using index access or taking the address/
		_ = item
	}
}
