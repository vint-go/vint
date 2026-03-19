package fixtures

func rangeExamples() {
	m := map[string]int{"a": 1, "b": 2}
	s := []string{"x", "y", "z"}

	// OK: using both key and value
	for k, v := range m {
		_ = k
		_ = v
	}

	// OK: using only key (single var form)
	for k := range m {
		_ = k
	}

	// Bad: key and blank value
	for k, _ := range m { // MATCH /should omit 2nd value from range; this loop is equivalent to `for k := range ...`/
		_ = k
	}

	// Bad: key and blank value on slice
	for i, _ := range s { // MATCH /should omit 2nd value from range; this loop is equivalent to `for i := range ...`/
		_ = i
	}

	// OK: using value
	for _, v := range s {
		_ = v
	}
}
