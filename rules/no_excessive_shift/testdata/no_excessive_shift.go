package fixtures

func excessiveShiftUint8() {
	var x uint8
	// Bad: shifting uint8 by 8 bits always yields 0
	y := x << 8 // MATCH /shift count 8 exceeds the bit width of type uint8 (8 bits)/
	_ = y
}

func excessiveShiftInt32() {
	var x int32
	// Bad: shifting int32 by 32 bits exceeds its width
	y := x << 32 // MATCH /shift count 32 exceeds the bit width of type int32 (32 bits)/
	_ = y
}

func excessiveShiftUint16() {
	var x uint16
	y := x >> 16 // MATCH /shift count 16 exceeds the bit width of type uint16 (16 bits)/
	_ = y
}

func excessiveShiftInt64() {
	var x int64
	y := x << 64 // MATCH /shift count 64 exceeds the bit width of type int64 (64 bits)/
	_ = y
}

func validShiftUint8() {
	var x uint8
	// Good: shift within the type's bit width
	y := x << 4
	_ = y
}

func validShiftInt64() {
	var x int64
	// Good: shift within the type's bit width
	y := x << 32
	_ = y
}

func validShiftInt32() {
	var x int32
	y := x >> 31
	_ = y
}

func validShiftUint16() {
	var x uint16
	y := x << 15
	_ = y
}
