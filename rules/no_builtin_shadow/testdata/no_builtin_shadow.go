package fixtures

import "fmt"

func invalidShadowLen() {
	len := 10 // MATCH /assignment shadows predeclared identifier len/
	_ = len
}

func invalidShadowError() {
	error := fmt.Errorf("something failed") // MATCH /assignment shadows predeclared identifier error/
	_ = error
}

func invalidShadowAppend() {
	append := func(s []int, v ...int) []int { return s } // MATCH /assignment shadows predeclared identifier append/
	_ = append
}

func invalidShadowMake() {
	make := 42 // MATCH /assignment shadows predeclared identifier make/
	_ = make
}

func invalidShadowNew() {
	new := "something" // MATCH /assignment shadows predeclared identifier new/
	_ = new
}

func invalidShadowCap() {
	cap := 100 // MATCH /assignment shadows predeclared identifier cap/
	_ = cap
}

func invalidShadowNil() {
	nil := 0 // MATCH /assignment shadows predeclared identifier nil/
	_ = nil
}

func invalidShadowTrue() {
	true := 1 // MATCH /assignment shadows predeclared identifier true/
	_ = true
}

func invalidShadowString() {
	string := "hello" // MATCH /assignment shadows predeclared identifier string/
	_ = string
}

func invalidShadowInt() {
	int := 5 // MATCH /assignment shadows predeclared identifier int/
	_ = int
}

func invalidShadowClose() {
	close := func() {} // MATCH /assignment shadows predeclared identifier close/
	_ = close
}

func invalidShadowRecover() {
	recover := func() interface{} { return nil } // MATCH /assignment shadows predeclared identifier recover/
	_ = recover
}

func invalidShadowPanic() {
	panic := func(v interface{}) {} // MATCH /assignment shadows predeclared identifier panic/
	_ = panic
}

// Valid examples - no matches expected
func validLength() {
	length := 10
	_ = length
}

func validErr() {
	err := fmt.Errorf("something failed")
	_ = err
}

func validCapacity() {
	capacity := 100
	_ = capacity
}

func validResult() {
	result := "hello"
	_ = result
}

func validNum() {
	num := 5
	_ = num
}
