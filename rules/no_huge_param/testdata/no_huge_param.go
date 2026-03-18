package fixtures

import "fmt"

// Large array type: [1024]int is 8192 bytes on amd64 (1024 * 8).
func processData(data [1024]int) { // MATCH /parameter 'data' exceeds the size threshold of 80 bytes with a size of 8192 bytes, consider passing it by pointer/
	_ = data
}

// Large struct type: 11 fields of 8 bytes each = 88 bytes > 80.
type LargeStruct struct {
	A, B, C, D, E int64
	F, G, H, I, J int64
	K             int64
}

func processLargeStruct(s LargeStruct) { // MATCH /parameter 's' exceeds the size threshold of 80 bytes with a size of 88 bytes, consider passing it by pointer/
	_ = s
}

// Valid: pointer to large array -- no copying.
func processDataPointer(data *[1024]int) {
	_ = data
}

// Valid: pointer to large struct.
func processLargeStructPointer(s *LargeStruct) {
	_ = s
}

// Valid: small struct under threshold.
type SmallStruct struct {
	X, Y int64
}

func processSmallStruct(s SmallStruct) {
	_ = s
}

// Valid: primitive types are small.
func processInt(x int) {
	_ = x
}

func processString(s string) {
	_ = s
}

// Valid: String() method is excluded (Stringer interface).
type MyStruct struct {
	A, B, C, D, E int64
	F, G, H, I, J int64
	K             int64
}

func (s MyStruct) String() string {
	return fmt.Sprintf("%v", s)
}

// Method that is NOT String() should check its receiver too.
func (s MyStruct) Process() { // MATCH /parameter 's' exceeds the size threshold of 80 bytes with a size of 88 bytes, consider passing it by pointer/
	_ = s
}

// Valid: pointer receiver on large struct.
func (s *MyStruct) Update() {
	s.A = 1
}
