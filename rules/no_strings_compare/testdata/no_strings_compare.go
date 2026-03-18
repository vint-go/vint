package fixtures

import "strings"

func invalid(a, b string) {
	if strings.Compare(a, b) == 0 { // MATCH /use == or < or > operators instead of strings.Compare/
		// equal
	}

	if strings.Compare(a, b) < 0 { // MATCH /use == or < or > operators instead of strings.Compare/
		// a < b
	}

	if strings.Compare(a, b) > 0 { // MATCH /use == or < or > operators instead of strings.Compare/
		// a > b
	}

	result := strings.Compare(a, b) // MATCH /use == or < or > operators instead of strings.Compare/
	_ = result
}

func valid(a, b string) {
	if a == b {
		// equal
	}

	if a < b {
		// a < b
	}

	if a > b {
		// a > b
	}

	// Other strings functions are fine
	_ = strings.Contains(a, b)
	_ = strings.HasPrefix(a, b)
}
