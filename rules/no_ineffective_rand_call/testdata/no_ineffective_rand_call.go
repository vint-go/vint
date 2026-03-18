package fixtures

import "math/rand"

func ineffectiveIntn() int {
	return rand.Intn(1) // MATCH /rand.Intn(1) always returns 0, this is likely a logic error/
}

func ineffectiveInt31n() int32 {
	return rand.Int31n(1) // MATCH /rand.Int31n(1) always returns 0, this is likely a logic error/
}

func ineffectiveInt63n() int64 {
	return rand.Int63n(1) // MATCH /rand.Int63n(1) always returns 0, this is likely a logic error/
}

// Valid cases - no match expected

func validIntn(max int) int {
	return rand.Intn(max) // dynamic argument
}

func validIntnWithTwo() int {
	return rand.Intn(2) // argument is not 1
}

func validIntnWithTen() int {
	return rand.Intn(10) // argument is not 1
}

func validInt31n(max int32) int32 {
	return rand.Int31n(max) // dynamic argument
}

func validRandInt() int {
	return rand.Int() // no argument
}
