package fixtures

func handleEmpty() {}
func handleNonEmpty() {}

func invalid(s string) {
	if len(s) == 0 { // MATCH /use direct string comparison instead of len() comparison/
		handleEmpty()
	}

	if len(s) != 0 { // MATCH /use direct string comparison instead of len() comparison/
		handleNonEmpty()
	}

	if 0 == len(s) { // MATCH /use direct string comparison instead of len() comparison/
		handleEmpty()
	}

	if 0 != len(s) { // MATCH /use direct string comparison instead of len() comparison/
		handleNonEmpty()
	}
}

func valid(s string) {
	if s == "" {
		handleEmpty()
	}

	if s != "" {
		handleNonEmpty()
	}

	// len() in non-comparison contexts is fine
	n := len(s)
	_ = n

	// len() with non-zero comparisons is fine
	if len(s) > 1 {
		handleNonEmpty()
	}

	if len(s) == 5 {
		handleNonEmpty()
	}

	if len(s) < 3 {
		handleEmpty()
	}

	// len() on non-string types should not be flagged
	sl := []int{1, 2, 3}
	if len(sl) == 0 {
		handleEmpty()
	}
	if 0 == len(sl) {
		handleEmpty()
	}

	m := map[string]int{"a": 1}
	if len(m) == 0 {
		handleEmpty()
	}
	if len(m) != 0 {
		handleNonEmpty()
	}

	arr := [3]int{1, 2, 3}
	if len(arr) == 0 {
		handleEmpty()
	}

	ch := make(chan int, 5)
	if len(ch) == 0 {
		handleEmpty()
	}
}
