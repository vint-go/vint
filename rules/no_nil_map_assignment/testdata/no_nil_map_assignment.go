package fixtures

func nilMapAssignmentBad() {
	var m map[string]int
	// Panic: assignment to entry in nil map
	m["key"] = 1 // MATCH /assignment to nil map/
}

func nilMapAssignmentGoodMake() {
	m := make(map[string]int)
	m["key"] = 1
}

func nilMapAssignmentGoodLiteral() {
	m := map[string]int{}
	m["key"] = 1
}

func nilMapAssignmentGoodInitialized() {
	var m map[string]int
	m = make(map[string]int)
	m["key"] = 1
}

func nilMapAssignmentGoodParam(m map[string]int) {
	m["key"] = 1
}

func nilMapAssignmentMultipleVars() {
	var a map[string]int
	var b map[string]string
	a["x"] = 1 // MATCH /assignment to nil map/
	b["y"] = "z" // MATCH /assignment to nil map/
}

func nilMapAssignmentGoodReassigned() {
	var m map[string]int
	m = getMap()
	m["key"] = 1
}

func getMap() map[string]int {
	return make(map[string]int)
}

func nilMapAssignmentGoodShortDecl() {
	m := map[string]int{"a": 1}
	m["b"] = 2
}
