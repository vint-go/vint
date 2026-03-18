package fixtures

func lookupBad(m map[string]int, key []byte) int {
	// Intermediate variable prevents compiler optimization
	s := string(key) // this assigns string([]byte) to s
	return m[s]      // MATCH /use string(b) directly as map key instead of intermediate variable to allow compiler optimization/
}

func lookupGood(m map[string]int, key []byte) int {
	// Direct conversion allows compiler optimization
	return m[string(key)]
}

func lookupGoodStringVar(m map[string]int) int {
	// Using a normal string variable is fine
	s := "hello"
	return m[s]
}

func lookupGoodStringLiteral(m map[string]int) int {
	// Using a string literal is fine
	return m["hello"]
}

func lookupBadWithAssign(m map[string]int, key []byte) int {
	var s string
	s = string(key)
	return m[s] // MATCH /use string(b) directly as map key instead of intermediate variable to allow compiler optimization/
}

func lookupGoodIntMap(m map[int]string, key int) string {
	// Not a string-keyed map, so no issue
	return m[key]
}

func lookupGoodUsedElsewhere(m map[string]int, key []byte) (int, string) {
	// If the converted string is used for other purposes, it might be intentional
	// but the rule still flags map key usage
	s := string(key)
	_ = len(s)
	return m[s], s // MATCH /use string(b) directly as map key instead of intermediate variable to allow compiler optimization/
}
