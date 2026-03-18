package fixtures

import "fmt"

// Invalid: converting to []rune before ranging over a string
func processRuneSlice(s string) {
	for _, r := range []rune(s) { // MATCH /range over string directly instead of converting to []rune/
		fmt.Printf("%c", r)
	}
}

// Invalid: same pattern with index and value
func processRuneSliceWithIndex(s string) {
	for i, r := range []rune(s) { // MATCH /range over string directly instead of converting to []rune/
		fmt.Printf("%d: %c", i, r)
	}
}

// Invalid: index-only range over []rune conversion
func processRuneSliceIndexOnly(s string) {
	for i := range []rune(s) { // MATCH /range over string directly instead of converting to []rune/
		fmt.Println(i)
	}
}

// Valid: range over string directly
func processStringDirect(s string) {
	for _, r := range s {
		fmt.Printf("%c", r)
	}
}

// Valid: range over an actual []rune variable
func processRuneVar() {
	runes := []rune("hello")
	for _, r := range runes {
		fmt.Printf("%c", r)
	}
}

// Valid: range over []byte conversion (different type)
func processByteSlice(s string) {
	for _, b := range []byte(s) {
		fmt.Printf("%x", b)
	}
}
