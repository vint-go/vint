package fixtures

import "math"

func uselessMathCeil(x int) float64 {
	return math.Ceil(float64(x)) // MATCH /math.Ceil called on a float64 converted from an integer, which is pointless/
}

func uselessMathFloor(x int) float64 {
	return math.Floor(float64(x)) // MATCH /math.Floor called on a float64 converted from an integer, which is pointless/
}

func uselessMathRound(x int) float64 {
	return math.Round(float64(x)) // MATCH /math.Round called on a float64 converted from an integer, which is pointless/
}

func uselessMathTrunc(x int) float64 {
	return math.Trunc(float64(x)) // MATCH /math.Trunc called on a float64 converted from an integer, which is pointless/
}

func uselessMathRoundToEven(x int) float64 {
	return math.RoundToEven(float64(x)) // MATCH /math.RoundToEven called on a float64 converted from an integer, which is pointless/
}

func uselessMathUnsigned(x uint) float64 {
	return math.Ceil(float64(x)) // MATCH /math.Ceil called on a float64 converted from an integer, which is pointless/
}

func uselessMathInt64(x int64) float64 {
	return math.Floor(float64(x)) // MATCH /math.Floor called on a float64 converted from an integer, which is pointless/
}

// Valid cases - no match expected

func validMathCeil(x float64) float64 {
	return math.Ceil(x) // no conversion from int
}

func validMathFloor(x float64) float64 {
	return math.Floor(x) // no conversion from int
}

func validMathRound(x float64) float64 {
	return math.Round(x) // no conversion from int
}

func validMathWithExpr(x float64) float64 {
	return math.Ceil(x * 1.5) // not a simple int-to-float conversion
}

func validDirectFloat64(x int) float64 {
	return float64(x) // just a conversion, no math call
}
