package fixtures

import (
	"encoding/binary"
	"math"
	"unsafe"
)

// Invalid: Using unsafe.Pointer to cast between types
func unsafeCast(i *int) *float64 {
	return (*float64)(unsafe.Pointer(i)) // MATCH /use of unsafe.Pointer bypasses Go type safety/
}

// Invalid: Using unsafe.Slice to create a slice from a pointer
func makeSlice(ptr *byte, length int) []byte {
	return unsafe.Slice(ptr, length) // MATCH /use of unsafe.Slice bypasses Go type safety/
}

// Invalid: Using unsafe.String to create a string from bytes
func makeString(ptr *byte, length int) string {
	return unsafe.String(ptr, length) // MATCH /use of unsafe.String bypasses Go type safety/
}

// Valid: Using encoding/binary for type conversion
func convertToFloat(data []byte) float64 {
	bits := binary.LittleEndian.Uint64(data)
	return math.Float64frombits(bits)
}

// Valid: Using standard library functions instead of unsafe
func copyBytes(src []byte) []byte {
	dst := make([]byte, len(src))
	copy(dst, src)
	return dst
}
