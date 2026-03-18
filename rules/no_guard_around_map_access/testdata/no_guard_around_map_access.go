package fixtures

func guardAroundMapAccessInvalid(m map[string]int, key string) int {
	var v int
	if _, ok := m[key]; ok { // MATCH /unnecessary guard around map access/
		v = m[key]
	}
	return v
}

func guardAroundMapAccessInvalidIntKey(m map[int]string, key int) string {
	var v string
	if _, ok := m[key]; ok { // MATCH /unnecessary guard around map access/
		v = m[key]
	}
	return v
}

func guardAroundMapAccessValid(m map[string]int, key string) int {
	return m[key]
}

func guardAroundMapAccessValidWithElse(m map[string]int, key string) int {
	var v int
	if _, ok := m[key]; ok {
		v = m[key]
	} else {
		v = -1
	}
	return v
}

func guardAroundMapAccessValidMultipleStatements(m map[string]int, key string) int {
	var v int
	if _, ok := m[key]; ok {
		v = m[key]
		println("found")
	}
	return v
}

func guardAroundMapAccessValidDifferentKey(m map[string]int, key string, other string) int {
	var v int
	if _, ok := m[key]; ok {
		v = m[other]
	}
	return v
}

func guardAroundMapAccessValidDifferentMap(m1 map[string]int, m2 map[string]int, key string) int {
	var v int
	if _, ok := m1[key]; ok {
		v = m2[key]
	}
	return v
}

func guardAroundMapAccessValidNotAssignment(m map[string]int, key string) {
	if _, ok := m[key]; ok {
		println(key)
	}
}

func guardAroundMapAccessValidUsesValue(m map[string]int, key string) int {
	if v, ok := m[key]; ok {
		return v
	}
	return 0
}

func guardAroundMapAccessValidDeleteInBody(m map[string]int, key string) {
	if _, ok := m[key]; ok {
		delete(m, key)
	}
}
