package fixtures

// Invalid: basic slide pattern shifting elements forward
func slideBasic() {
	s := make([]int, 10)
	for i := 0; i < len(s)-1; i++ { // MATCH /should replace loop with copy(s, s[1:])/
		s[i] = s[i+1]
	}
}

// Valid: uses copy directly
func slideWithCopy() {
	s := make([]int, 10)
	copy(s, s[1:])
}

// Valid: loop body has more than one statement
func slideWithExtra() {
	s := make([]int, 10)
	for i := 0; i < len(s)-1; i++ {
		s[i] = s[i+1]
		_ = i
	}
}

// Valid: different slices on left and right
func slideDifferentSlices() {
	s := make([]int, 10)
	t := make([]int, 10)
	for i := 0; i < len(s)-1; i++ {
		t[i] = s[i+1]
	}
}

// Valid: not shifting by 1
func slideByTwo() {
	s := make([]int, 10)
	for i := 0; i < len(s)-1; i++ {
		s[i] = s[i+2]
	}
}

// Valid: not starting from 0
func slideFromOne() {
	s := make([]int, 10)
	for i := 1; i < len(s)-1; i++ {
		s[i] = s[i+1]
	}
}

// Valid: using <= instead of <
func slideWrongOperator() {
	s := make([]int, 10)
	for i := 0; i <= len(s)-1; i++ {
		s[i] = s[i+1]
	}
}

// Valid: decrementing instead of incrementing
func slideDecrement() {
	s := make([]int, 10)
	for i := 0; i < len(s)-1; i-- {
		s[i] = s[i+1]
	}
}
