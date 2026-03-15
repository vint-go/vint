package fixtures

import (
	"fmt"
	"strconv"
)

func badStringConversion() {
	n := 65
	s := string(n) // MATCH /string(int) produces a rune character, not a decimal string; use strconv.Itoa or fmt.Sprint instead/
	_ = s
}

func badReturnStringConversion(code int) string {
	return string(code) // MATCH /string(int) produces a rune character, not a decimal string; use strconv.Itoa or fmt.Sprint instead/
}

func goodStrconvItoa() {
	n := 65
	s := strconv.Itoa(n)
	_ = s
}

func goodFmtSprint(code int) string {
	return fmt.Sprint(code)
}

func goodRuneConversion() {
	s := string(rune(65))
	_ = s
}

func goodByteSlice() {
	s := string([]byte{65})
	_ = s
}

func goodStringLiteral() {
	s := string("hello")
	_ = s
}
