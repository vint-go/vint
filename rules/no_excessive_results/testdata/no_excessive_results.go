package fixtures

// Too many return values (6 > 3)
func getDetails() (string, int, bool, error, []byte, map[string]string) { // MATCH /function getDetails has too many results (6 > 3), consider grouping into a struct/
	return "", 0, false, nil, nil, nil
}

// Exactly at the limit (3 results) - should not trigger
func atLimit() (string, int, bool) {
	return "", 0, false
}

// Below the limit (2 results) - should not trigger
func belowLimit() (string, error) {
	return "", nil
}

// No return values - should not trigger
func noResults() {
}

// One return value - should not trigger
func oneResult() error {
	return nil
}

// Four results - over the limit of 3
func fourResults() (int, string, bool, error) { // MATCH /function fourResults has too many results (4 > 3), consider grouping into a struct/
	return 0, "", false, nil
}

type MyStruct struct{}

// Method with too many results
func (m *MyStruct) tooManyResults() (int, string, bool, float64) { // MATCH /function (*MyStruct).tooManyResults has too many results (4 > 3), consider grouping into a struct/
	return 0, "", false, 0.0
}

// Method with acceptable results
func (m *MyStruct) okResults() (int, error) {
	return 0, nil
}

// Valid: grouped into a struct
type Details struct {
	Name   string
	Count  int
	Active bool
	Data   []byte
	Meta   map[string]string
}

func getDetailsGrouped() (Details, error) {
	return Details{}, nil
}
