package fixtures

import "unicode/utf8"

func badDecodeRune() {
	s := "hello"
	_ = []rune(s)[0] // MATCH /use utf8.DecodeRuneInString instead of []rune(s)[0] to avoid full string-to-rune-slice conversion/
}

func badDecodeRuneInline() {
	_ = []rune("world")[0] // MATCH /use utf8.DecodeRuneInString instead of []rune(s)[0] to avoid full string-to-rune-slice conversion/
}

func goodDecodeRune() {
	s := "hello"
	r, _ := utf8.DecodeRuneInString(s)
	_ = r
}

func goodRuneSliceFullAccess() {
	s := "hello"
	runes := []rune(s)
	_ = runes[0]
}

func goodRuneSliceNonZeroIndex() {
	s := "hello"
	_ = []rune(s)[1]
}

func goodByteSliceAccess() {
	s := "hello"
	_ = []byte(s)[0]
}
