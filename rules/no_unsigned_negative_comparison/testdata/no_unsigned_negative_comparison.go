package fixtures

func unsignedLessThanZero(x uint) bool {
	return x < 0 // MATCH /comparing unsigned value against negative value is always false/
}

func unsignedGreaterEqualZero(x uint) bool {
	return x >= 0 // MATCH /comparing unsigned value against negative value is always true/
}

func unsignedLessThanNegOne(x uint) bool {
	return x < -1 // MATCH /comparing unsigned value against negative value is always false/
}

func unsignedGreaterThanNegOne(x uint) bool {
	return x > -1 // MATCH /comparing unsigned value against negative value is always true/
}

func negativeOnLeft(x uint) bool {
	return -1 > x // MATCH /comparing unsigned value against negative value is always false/
}

func negativeOnLeftGEQ(x uint) bool {
	return -1 <= x // MATCH /comparing unsigned value against negative value is always true/
}

func uint8Type(x uint8) bool {
	return x < 0 // MATCH /comparing unsigned value against negative value is always false/
}

func uint16Type(x uint16) bool {
	return x < 0 // MATCH /comparing unsigned value against negative value is always false/
}

func uint32Type(x uint32) bool {
	return x < 0 // MATCH /comparing unsigned value against negative value is always false/
}

func uint64Type(x uint64) bool {
	return x < 0 // MATCH /comparing unsigned value against negative value is always false/
}

func uintptrType(x uintptr) bool {
	return x < 0 // MATCH /comparing unsigned value against negative value is always false/
}

// Valid cases - should NOT trigger

func signedComparison(x int) bool {
	return x < 0
}

func unsignedEqualZero(x uint) bool {
	return x == 0
}

func unsignedNotEqualZero(x uint) bool {
	return x != 0
}

func unsignedGreaterThanPositive(x uint) bool {
	return x > 5
}

func unsignedLessThanPositive(x uint) bool {
	return x < 10
}
