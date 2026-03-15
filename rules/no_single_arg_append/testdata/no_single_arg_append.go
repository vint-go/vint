package fixtures

func singleArgAppendInvalid() {
	s := []int{1, 2, 3}
	s = append(s) // MATCH /append with only one argument has no effect/
}

func singleArgAppendValid() {
	s := []int{1, 2, 3}
	s = append(s, 4)
}

func singleArgAppendValidVariadic() {
	s := []int{1, 2, 3}
	t := []int{4, 5}
	s = append(s, t...)
}

func singleArgAppendValidMultiple() {
	s := []int{1, 2, 3}
	s = append(s, 4, 5, 6)
}
