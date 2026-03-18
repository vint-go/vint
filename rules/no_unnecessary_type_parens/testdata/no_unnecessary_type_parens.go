package fixtures

// Invalid: pointer with parenthesized base type
var x *(int) // MATCH /unnecessary parentheses in type expression/

// Invalid: map with parenthesized key type
var m1 map[(string)]int // MATCH /unnecessary parentheses in type expression/

// Invalid: map with parenthesized value type
var m2 map[string](int) // MATCH /unnecessary parentheses in type expression/

// Invalid: array with parenthesized element type
var a [3](int) // MATCH /unnecessary parentheses in type expression/

// Invalid: slice with parenthesized element type
var s [](int) // MATCH /unnecessary parentheses in type expression/

// Invalid: channel with parenthesized value type
var ch chan (int) // MATCH /unnecessary parentheses in type expression/

// Invalid: function parameter with parenthesized type
func badParam(x (int)) {} // MATCH /unnecessary parentheses in type expression/

// Invalid: type assertion with parenthesized type
func badAssert(x interface{}) {
	_ = x.((int)) // MATCH /unnecessary parentheses in type expression/
}

// Valid: pointer without unnecessary parens
var x2 *int

// Valid: map without unnecessary parens
var m3 map[string]int

// Valid: array without unnecessary parens
var a2 [3]int

// Valid: slice without unnecessary parens
var s2 []int

// Valid: channel without unnecessary parens
var ch2 chan int

// Valid: function parameter without unnecessary parens
func goodParam(x int) {}

// Valid: function return without unnecessary parens
func goodReturn() int {
	return 0
}

// Valid: type assertion without unnecessary parens
func goodAssert(x interface{}) {
	_ = x.(int)
}

// Valid: multiple return values in parens are fine (these are required)
func multiReturn() (int, error) {
	return 0, nil
}
