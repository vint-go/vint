package fixtures

func badCompositeLiteralNilCheck() {
	m := map[string]int{}
	if m == nil { // MATCH /checking never-nil value m against nil/
		panic("impossible")
	}
}

func badSliceLiteralNilCheck() {
	s := []int{1, 2, 3}
	if s == nil { // MATCH /checking never-nil value s against nil/
		panic("impossible")
	}
}

func badStructLiteralNilCheck() {
	type MyStruct struct{ X int }
	s := &MyStruct{}
	if s == nil { // MATCH /checking never-nil value s against nil/
		panic("impossible")
	}
}

func badDirectCompositeLitNilCheck() {
	if (map[string]int{}) == nil { // MATCH /checking never-nil value (map[string]int{}) against nil/
		panic("impossible")
	}
}

func badNewNilCheck() {
	p := new(int)
	if p == nil { // MATCH /checking never-nil value p against nil/
		panic("impossible")
	}
}

func badMakeNilCheck() {
	ch := make(chan int)
	if ch == nil { // MATCH /checking never-nil value ch against nil/
		panic("impossible")
	}
}

func badNilOnLeftSide() {
	m := map[string]int{}
	if nil == m { // MATCH /checking never-nil value m against nil/
		panic("impossible")
	}
}

func badNotEqual() {
	m := map[string]int{}
	if m != nil { // MATCH /checking never-nil value m against nil/
		// always true
	}
}

func goodParameterNilCheck(m map[string]int) {
	// m is a parameter, could be nil
	if m == nil {
		m = make(map[string]int)
	}
}

func goodReassigned() {
	m := map[string]int{}
	m = getMap()
	if m == nil {
		// m was reassigned, could be nil
	}
}

func goodNilCheckOnFuncResult() {
	m := getMap()
	if m == nil {
		// function result could be nil
	}
}

func goodNilCheckOnVariable() {
	var m map[string]int
	if m == nil {
		m = make(map[string]int)
	}
	_ = m
}

func getMap() map[string]int {
	return nil
}
