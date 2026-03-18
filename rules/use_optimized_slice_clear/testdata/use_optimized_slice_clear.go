package fixtures

type myStruct struct {
	x int
	y string
}

// Invalid: classic for loop clearing int slice with 0
func clearIntSlice() {
	s := make([]int, 10)
	for i := 0; i < len(s); i++ { // MATCH /use 'for i := range s' to enable compiler memclr optimization/
		s[i] = 0
	}
}

// Invalid: classic for loop clearing string slice with ""
func clearStringSlice() {
	s := make([]string, 10)
	for i := 0; i < len(s); i++ { // MATCH /use 'for i := range s' to enable compiler memclr optimization/
		s[i] = ""
	}
}

// Invalid: classic for loop clearing bool slice with false
func clearBoolSlice() {
	s := make([]bool, 10)
	for i := 0; i < len(s); i++ { // MATCH /use 'for i := range s' to enable compiler memclr optimization/
		s[i] = false
	}
}

// Invalid: classic for loop clearing pointer slice with nil
func clearPointerSlice() {
	s := make([]*int, 10)
	for i := 0; i < len(s); i++ { // MATCH /use 'for i := range s' to enable compiler memclr optimization/
		s[i] = nil
	}
}

// Invalid: classic for loop clearing struct slice with zero literal
func clearStructSlice() {
	s := make([]myStruct, 10)
	for i := 0; i < len(s); i++ { // MATCH /use 'for i := range s' to enable compiler memclr optimization/
		s[i] = myStruct{}
	}
}

// Valid: already uses range (optimized idiom)
func clearWithRange() {
	s := make([]int, 10)
	for i := range s {
		s[i] = 0
	}
}

// Valid: body has more than one statement
func clearWithExtraWork() {
	s := make([]int, 10)
	sum := 0
	for i := 0; i < len(s); i++ {
		sum += s[i]
		s[i] = 0
	}
	_ = sum
}

// Valid: not assigning zero value
func setNonZero() {
	s := make([]int, 10)
	for i := 0; i < len(s); i++ {
		s[i] = 1
	}
}

// Valid: using different slice in index expression
func differentSlice() {
	s := make([]int, 10)
	t := make([]int, 10)
	for i := 0; i < len(s); i++ {
		t[i] = 0
	}
}

// Valid: not starting from 0
func notStartingFromZero() {
	s := make([]int, 10)
	for i := 1; i < len(s); i++ {
		s[i] = 0
	}
}

// Valid: using different condition variable
func differentCondVar() {
	s := make([]int, 10)
	n := len(s)
	for i := 0; i < n; i++ {
		s[i] = 0
	}
}
