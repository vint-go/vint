package fixtures

// Invalid: two same unnamed types (ambiguous, neither is error/bool)
func twoSameFloat() (float64, float64) { // MATCH /unnamed results of type 'float64', 'float64' may benefit from named results/
	return 0.0, 0.0
}

// Invalid: two same unnamed int types
func twoSameInt() (int, int) { // MATCH /unnamed results of type 'int', 'int' may benefit from named results/
	return 0, 0
}

// Invalid: more than 2 results with duplicate types (excluding trailing error)
func threeWithDup() (int, int, error) { // MATCH /unnamed results of type 'int', 'int', 'error' may benefit from named results/
	return 0, 0, nil
}

// Invalid: more than 2 results with duplicate types (excluding trailing bool)
func threeWithDupBool() (string, string, bool) { // MATCH /unnamed results of type 'string', 'string', 'bool' may benefit from named results/
	return "", "", false
}

// Invalid: four results with duplicate types
func fourWithDup() (int, string, int, error) { // MATCH /unnamed results of type 'int', 'string', 'int', 'error' may benefit from named results/
	return 0, "", 0, nil
}

// Valid: named results
func namedResults() (x, y float64) {
	return 0.0, 0.0
}

// Valid: single return
func singleReturn() int {
	return 0
}

// Valid: int and error (different types for 2-return)
func intAndError() (int, error) {
	return 0, nil
}

// Valid: int and bool (different types for 2-return)
func intAndBool() (int, bool) {
	return 0, false
}

// Valid: three distinct types with trailing error (no duplicates after excluding error)
func threeDistinctWithError() (int, string, error) {
	return 0, "", nil
}

// Valid: no return values
func noReturn() {
}

// Valid: two same type error (second is error, excluded)
func twoErrors() (error, error) {
	return nil, nil
}

// Valid: two same type bool (second is bool, excluded)
func twoBools() (bool, bool) {
	return false, false
}

// Valid: string and int (different types for 2-return)
func stringAndInt() (string, int) {
	return "", 0
}
