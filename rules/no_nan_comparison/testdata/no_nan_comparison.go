package fixtures

import "math"

// Invalid: comparing with math.NaN() using ==
func badEqual(x float64) bool {
	return x == math.NaN() // MATCH /comparing with math.NaN() is always false, use math.IsNaN() instead/
}

// Invalid: comparing with math.NaN() using !=
func badNotEqual(x float64) bool {
	return x != math.NaN() // MATCH /comparing with math.NaN() is always false, use math.IsNaN() instead/
}

// Invalid: comparing with math.NaN() using <
func badLessThan(x float64) bool {
	return x < math.NaN() // MATCH /comparing with math.NaN() is always false, use math.IsNaN() instead/
}

// Invalid: comparing with math.NaN() using >
func badGreaterThan(x float64) bool {
	return x > math.NaN() // MATCH /comparing with math.NaN() is always false, use math.IsNaN() instead/
}

// Invalid: comparing with math.NaN() on the left side
func badNaNOnLeft(x float64) bool {
	return math.NaN() == x // MATCH /comparing with math.NaN() is always false, use math.IsNaN() instead/
}

// Valid: using math.IsNaN() instead
func goodIsNaN(x float64) bool {
	return math.IsNaN(x)
}

// Valid: comparing two non-NaN values
func goodRegularComparison(x, y float64) bool {
	return x == y
}

// Valid: calling math.NaN() without comparison
func goodNaNAssign() float64 {
	return math.NaN()
}
