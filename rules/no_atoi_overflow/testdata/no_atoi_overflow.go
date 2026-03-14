package fixtures

import (
	"log"
	"math"
	"strconv"
)

// Invalid: direct conversion of Atoi result to int32 without bounds check
func noAtoiOverflowDirectInt32(input string) int32 {
	val, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
	}
	smallVal := int32(val) // MATCH /potential integer overflow: conversion of strconv.Atoi result to int32 without bounds check/
	return smallVal
}

// Invalid: direct conversion to int16
func noAtoiOverflowDirectInt16(input string) int16 {
	val, _ := strconv.Atoi(input)
	result := int16(val) // MATCH /potential integer overflow: conversion of strconv.Atoi result to int16 without bounds check/
	return result
}

// Invalid: direct conversion to int8
func noAtoiOverflowDirectInt8(input string) int8 {
	val, _ := strconv.Atoi(input)
	result := int8(val) // MATCH /potential integer overflow: conversion of strconv.Atoi result to int8 without bounds check/
	return result
}

// Invalid: direct conversion to uint32
func noAtoiOverflowDirectUint32(input string) uint32 {
	val, _ := strconv.Atoi(input)
	result := uint32(val) // MATCH /potential integer overflow: conversion of strconv.Atoi result to uint32 without bounds check/
	return result
}

// Invalid: direct conversion to uint16
func noAtoiOverflowDirectUint16(input string) uint16 {
	val, _ := strconv.Atoi(input)
	result := uint16(val) // MATCH /potential integer overflow: conversion of strconv.Atoi result to uint16 without bounds check/
	return result
}

// Invalid: direct conversion to uint8
func noAtoiOverflowDirectUint8(input string) uint8 {
	val, _ := strconv.Atoi(input)
	result := uint8(val) // MATCH /potential integer overflow: conversion of strconv.Atoi result to uint8 without bounds check/
	return result
}

// Valid: bounds checking before conversion
func noAtoiOverflowBoundsCheck(input string) int32 {
	val, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
	}
	if val > math.MaxInt32 || val < math.MinInt32 {
		log.Fatal("value out of range for int32")
	}
	smallVal := int32(val)
	return smallVal
}

// Valid: using strconv.ParseInt with explicit bit size
func noAtoiOverflowParseInt(input string) int32 {
	val, err := strconv.ParseInt(input, 10, 32)
	if err != nil {
		log.Fatal(err)
	}
	smallVal := int32(val)
	return smallVal
}

// Valid: converting to int64 (same or larger size)
func noAtoiOverflowInt64(input string) int64 {
	val, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
	}
	bigVal := int64(val)
	return bigVal
}

// Valid: no conversion at all
func noAtoiOverflowNoConversion(input string) int {
	val, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
	}
	return val
}

// Valid: blank identifier for Atoi result
func noAtoiOverflowBlankIdentifier(input string) {
	_, _ = strconv.Atoi(input)
}
