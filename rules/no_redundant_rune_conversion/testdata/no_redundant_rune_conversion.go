package fixtures

import "fmt"

// Invalid: unnecessary conversion to []rune before range loop
func processRedundant(s string) {
	for _, r := range []rune(s) { // MATCH /unnecessary conversion to []rune before range loop/
		fmt.Printf("%c", r)
	}
}

// Invalid: same pattern with index and value
func processRedundantWithIndex(s string) {
	for i, r := range []rune(s) { // MATCH /unnecessary conversion to []rune before range loop/
		fmt.Printf("%d: %c", i, r)
	}
}

// Invalid: with only index
func processRedundantIndexOnly(s string) {
	for i := range []rune(s) { // MATCH /unnecessary conversion to []rune before range loop/
		fmt.Println(i)
	}
}

// Valid: ranging directly over a string
func processDirect(s string) {
	for _, r := range s {
		fmt.Printf("%c", r)
	}
}

// Valid: converting to []rune and storing (not in range)
func convertAndStore(s string) {
	runes := []rune(s)
	for _, r := range runes {
		fmt.Printf("%c", r)
	}
}

// Valid: converting a []rune variable (not a string)
func processRuneSlice(rs []rune) {
	for _, r := range []rune(rs) {
		fmt.Printf("%c", r)
	}
}

// Valid: converting to []byte (not []rune)
func processBytes(s string) {
	for _, b := range []byte(s) {
		fmt.Printf("%d", b)
	}
}
