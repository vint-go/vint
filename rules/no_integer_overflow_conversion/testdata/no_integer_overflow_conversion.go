package fixtures

import (
	"fmt"
	"math"
)

// Invalid: converting int64 to int32 without bounds check
func intOverflowInt64ToInt32() {
	var bigVal int64 = math.MaxInt64
	smallVal := int32(bigVal) // MATCH /potential integer overflow: conversion of int64 to int32 may cause overflow or sign change/
	fmt.Println(smallVal)
}

// Invalid: converting uint to int may overflow on large values
func intOverflowUintToInt(size uint) {
	intSize := int(size) // MATCH /potential integer overflow: conversion of uint to int may cause overflow or sign change/
	fmt.Println(intSize)
}

// Invalid: converting signed to unsigned without checking for negative
func intOverflowIntToUint(val int) uint {
	return uint(val) // MATCH /potential integer overflow: conversion of int to uint may cause overflow or sign change/
}

// Invalid: converting int64 to int16
func intOverflowInt64ToInt16(val int64) int16 {
	return int16(val) // MATCH /potential integer overflow: conversion of int64 to int16 may cause overflow or sign change/
}

// Invalid: converting int64 to int8
func intOverflowInt64ToInt8(val int64) int8 {
	return int8(val) // MATCH /potential integer overflow: conversion of int64 to int8 may cause overflow or sign change/
}

// Invalid: converting uint64 to uint32
func intOverflowUint64ToUint32(val uint64) uint32 {
	return uint32(val) // MATCH /potential integer overflow: conversion of uint64 to uint32 may cause overflow or sign change/
}

// Invalid: converting int to int32
func intOverflowIntToInt32(val int) int32 {
	return int32(val) // MATCH /potential integer overflow: conversion of int to int32 may cause overflow or sign change/
}

// Valid: bounds checking before conversion
func intOverflowValidBoundsCheck() {
	var bigVal int64 = 42
	if bigVal > math.MaxInt32 || bigVal < math.MinInt32 {
		fmt.Println("value out of range")
		return
	}
	smallVal := int32(bigVal)
	fmt.Println(smallVal)
}

// Valid: safe conversion from smaller to larger type
func intOverflowValidSmallerToLarger(val int32) int64 {
	return int64(val)
}

// Valid: same type conversion
func intOverflowValidSameType(val int32) int32 {
	return int32(val)
}

// Valid: converting uint8 to int (always safe - smaller unsigned to larger signed)
func intOverflowValidUint8ToInt(val uint8) int {
	return int(val)
}

// Valid: converting int32 to int64 (always safe)
func intOverflowValidInt32ToInt64(val int32) int64 {
	return int64(val)
}

// Valid: converting uint16 to uint64 (always safe)
func intOverflowValidUint16ToUint64(val uint16) uint64 {
	return uint64(val)
}
