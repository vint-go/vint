package fixtures

// Invalid: casting x (int32) to smaller int16 before comparison truncates.
func truncatingInt32ToInt16(x int32, y int16) {
	if int16(x) < y { // MATCH /truncating conversion from int32 to int16 before comparison may lose information/
		_ = x
	}
}

// Invalid: casting int64 to int32 before comparison.
func truncatingInt64ToInt32(x int64, y int32) {
	if int32(x) > y { // MATCH /truncating conversion from int64 to int32 before comparison may lose information/
		_ = x
	}
}

// Invalid: casting int32 to int8 before comparison.
func truncatingInt32ToInt8(x int32, y int8) {
	if int8(x) == y { // MATCH /truncating conversion from int32 to int8 before comparison may lose information/
		_ = x
	}
}

// Invalid: casting uint32 to uint16 before comparison.
func truncatingUint32ToUint16(x uint32, y uint16) {
	if uint16(x) <= y { // MATCH /truncating conversion from uint32 to uint16 before comparison may lose information/
		_ = x
	}
}

// Invalid: casting int64 to int16 before comparison.
func truncatingInt64ToInt16(x int64, y int16) {
	if int16(x) != y { // MATCH /truncating conversion from int64 to int16 before comparison may lose information/
		_ = x
	}
}

// Invalid: casting uint64 to uint8 before comparison.
func truncatingUint64ToUint8(x uint64, y uint8) {
	if uint8(x) >= y { // MATCH /truncating conversion from uint64 to uint8 before comparison may lose information/
		_ = x
	}
}

// Invalid: the truncating conversion can be on the right side too.
func truncatingOnRight(x int16, y int32) {
	if x < int16(y) { // MATCH /truncating conversion from int32 to int16 before comparison may lose information/
		_ = x
	}
}

// Valid: casting the narrower operand to the larger type (widening).
func wideningComparison(x int32, y int16) {
	if x < int32(y) {
		_ = x
	}
}

// Valid: same-size conversion (int32 to int32 is identity).
func sameSizeConversion(x int32, y int32) {
	if int32(x) < y {
		_ = x
	}
}

// Valid: no conversion at all.
func noConversion(x int32, y int32) {
	if x < y {
		_ = x
	}
}

// Valid: conversion to a larger type.
func wideningInt16ToInt32(x int16, y int32) {
	if int32(x) < y {
		_ = x
	}
}

// Valid: architecture-dependent types are skipped by default.
func archDependentSkipped(x int, y int32) {
	if int32(x) < y {
		_ = x
	}
}
