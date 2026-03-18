package fixtures

func dubiousShiftUint8() {
	var x uint8 = 1
	// Shifting uint8 by 8 bits always yields 0
	y := x << 8 // MATCH /dubious bit shift of uint8 value by 8 bits, the type has only 8 bits/
	_ = y
}

func dubiousShiftUint8More() {
	var x uint8 = 1
	y := x << 10 // MATCH /dubious bit shift of uint8 value by 10 bits, the type has only 8 bits/
	_ = y
}

func dubiousShiftInt16() {
	var x int16 = 1
	y := x << 16 // MATCH /dubious bit shift of int16 value by 16 bits, the type has only 16 bits/
	_ = y
}

func dubiousShiftUint32() {
	var x uint32 = 1
	y := x << 32 // MATCH /dubious bit shift of uint32 value by 32 bits, the type has only 32 bits/
	_ = y
}

func dubiousShiftInt64() {
	var x int64 = 1
	y := x << 64 // MATCH /dubious bit shift of int64 value by 64 bits, the type has only 64 bits/
	_ = y
}

func dubiousRightShiftUint8() {
	var x uint8 = 255
	y := x >> 8 // MATCH /dubious bit shift of uint8 value by 8 bits, the type has only 8 bits/
	_ = y
}

func dubiousShiftInt32() {
	var x int32 = 1
	y := x >> 32 // MATCH /dubious bit shift of int32 value by 32 bits, the type has only 32 bits/
	_ = y
}

// Valid: shift within the type's bit width
func validShiftUint8() {
	var x uint8 = 1
	y := x << 4
	_ = y
}

func validShiftUint16() {
	var x uint16 = 1
	y := x << 8
	_ = y
}

func validShiftInt32() {
	var x int32 = 1
	y := x >> 31
	_ = y
}

func validShiftInt64() {
	var x int64 = 1
	y := x << 32
	_ = y
}

// Valid: platform-dependent sizes are not flagged
func validShiftInt() {
	var x int = 1
	y := x << 32
	_ = y
}

func validShiftUint() {
	var x uint = 1
	y := x << 32
	_ = y
}
