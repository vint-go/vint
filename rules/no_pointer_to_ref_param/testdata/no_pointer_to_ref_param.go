package fixtures

// Invalid: pointer to map in input parameter
func f1(m *map[string]int) { // MATCH /input parameter should not be a pointer to map/
}

// Invalid: pointer to channel in output parameter
func f2() *chan int { // MATCH /output parameter should not be a pointer to chan/
	return nil
}

// Invalid: pointer to map input and pointer to channel output together
func f3a(m *map[string]int) { // MATCH /input parameter should not be a pointer to map/
}

func f3b() *chan *int { // MATCH /output parameter should not be a pointer to chan/
	return nil
}

// Invalid: pointer to interface in input parameter
func f4(i *interface{}) { // MATCH /input parameter should not be a pointer to interface/
}

// Invalid: pointer to map as unnamed output
func f5() *map[string]string { // MATCH /output parameter should not be a pointer to map/
	return nil
}

// Invalid: pointer to channel in named output parameter
func f6() (ch *chan bool) { // MATCH /output parameter should not be a pointer to chan/
	return nil
}

// Valid: map used directly without pointer
func g1(m map[string]int) {
}

// Valid: channel used directly without pointer
func g2() chan *int {
	return nil
}

// Valid: interface used directly without pointer
func g3(i interface{}) {
}

// Valid: pointer to a struct is fine
type MyStruct struct{}

func g4(s *MyStruct) {
}

// Valid: pointer to a slice is fine
func g5(s *[]int) {
}

// Valid: pointer to a basic type is fine
func g6(n *int) {
}

// Valid: no parameters or return values
func g7() {
}
