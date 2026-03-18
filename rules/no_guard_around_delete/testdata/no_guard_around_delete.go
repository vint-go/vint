package fixtures

func guardAroundDeleteInvalid(m map[string]int, key string) {
	if _, ok := m[key]; ok { // MATCH /unnecessary guard around call to delete/
		delete(m, key)
	}
}

func guardAroundDeleteInvalidIntKey(m map[int]string, key int) {
	if _, ok := m[key]; ok { // MATCH /unnecessary guard around call to delete/
		delete(m, key)
	}
}

func guardAroundDeleteValid(m map[string]int, key string) {
	delete(m, key)
}

func guardAroundDeleteValidWithElse(m map[string]int, key string) {
	if _, ok := m[key]; ok {
		delete(m, key)
	} else {
		m[key] = 0
	}
}

func guardAroundDeleteValidMultipleStatements(m map[string]int, key string) {
	if _, ok := m[key]; ok {
		delete(m, key)
		println("deleted")
	}
}

func guardAroundDeleteValidDifferentKey(m map[string]int, key string, other string) {
	if _, ok := m[key]; ok {
		delete(m, other)
	}
}

func guardAroundDeleteValidDifferentMap(m1 map[string]int, m2 map[string]int, key string) {
	if _, ok := m1[key]; ok {
		delete(m2, key)
	}
}

func guardAroundDeleteValidNotDelete(m map[string]int, key string) {
	if _, ok := m[key]; ok {
		println(key)
	}
}
