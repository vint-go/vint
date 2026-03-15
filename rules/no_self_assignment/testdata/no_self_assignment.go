package fixtures

type MyStruct struct {
	Field int
}

func selfAssignSimple() {
	x := 5
	x = x // MATCH /useless self-assignment of x/
}

func selfAssignStructField(s *MyStruct) {
	s.Field = s.Field // MATCH /useless self-assignment of s.Field/
}

func validAssignDifferentVar() {
	x := 5
	y := x // assignment to a different variable
	_ = y
}

func validAssignDifferentValue(s *MyStruct) {
	s.Field = computeNewValue() // assignment with a different value
}

func validShortVarDecl() {
	x := 5
	_ = x
}

func validOpAssign() {
	x := 5
	x += 1
	_ = x
}

func validMultiAssign() {
	x, y := 1, 2
	x, y = y, x // swap is not self-assignment
	_, _ = x, y
}

func computeNewValue() int {
	return 42
}
