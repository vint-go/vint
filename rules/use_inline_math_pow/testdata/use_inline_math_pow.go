package fixtures

import "math"

func squareExample(x float64) float64 {
	return math.Pow(x, 2) // MATCH /math.Pow can be replaced by inline multiplication: x * x/
}

func cubeExample(x float64) float64 {
	return math.Pow(x, 3) // MATCH /math.Pow can be replaced by inline multiplication: x * x * x/
}

func fourthPowerExample(x float64) float64 {
	return math.Pow(x, 4) // MATCH /math.Pow can be replaced by inline multiplication: x * x * x * x/
}

func identityExample(x float64) float64 {
	return math.Pow(x, 1) // MATCH /math.Pow can be replaced by inline multiplication: x/
}

func zeroPowerExample(x float64) float64 {
	return math.Pow(x, 0) // MATCH /math.Pow can be replaced by inline multiplication: 1/
}

func floatExponentWholeNumber(x float64) float64 {
	return math.Pow(x, 2.0) // MATCH /math.Pow can be replaced by inline multiplication: x * x/
}

// Valid examples - should NOT trigger

func largePower(x float64) float64 {
	return math.Pow(x, 5) // exponent too large
}

func variablePower(x, n float64) float64 {
	return math.Pow(x, n) // not a literal
}

func fractionalPower(x float64) float64 {
	return math.Pow(x, 0.5) // not an integer exponent
}

func negativePower(x float64) float64 {
	return math.Pow(x, -1) // negative exponent
}

func validSquare(x float64) float64 {
	return x * x // already inline
}
