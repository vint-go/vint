package fixtures

// Invalid: mixed case hex digit letters.

func mixedCase() {
	x := 0xfF  // MATCH /hex literal 0xfF has mixed case letter digits, use consistent casing/
	y := 0xaB  // MATCH /hex literal 0xaB has mixed case letter digits, use consistent casing/
	z := 0xAbCd // MATCH /hex literal 0xAbCd has mixed case letter digits, use consistent casing/
	_ = x
	_ = y
	_ = z
}

// Invalid: uppercase 0X prefix.

func uppercasePrefix() {
	x := 0XFF // MATCH /hex literal 0XFF uses uppercase 0X prefix, use lowercase 0x instead/
	y := 0Xff // MATCH /hex literal 0Xff uses uppercase 0X prefix, use lowercase 0x instead/
	z := 0X00 // MATCH /hex literal 0X00 uses uppercase 0X prefix, use lowercase 0x instead/
	_ = x
	_ = y
	_ = z
}

// Valid: consistent lowercase hex digits.

func lowerCase() {
	x := 0xff
	y := 0xab
	z := 0xabcd
	_ = x
	_ = y
	_ = z
}

// Valid: consistent uppercase hex digits.

func upperCase() {
	x := 0xFF
	y := 0xAB
	z := 0xABCD
	_ = x
	_ = y
	_ = z
}

// Valid: no letter digits (only numeric digits).

func numericOnly() {
	x := 0x00
	y := 0x123
	z := 0x99
	_ = x
	_ = y
	_ = z
}

// Valid: non-hex integer literals.

func nonHex() {
	x := 42
	y := 0o77
	z := 0b1010
	_ = x
	_ = y
	_ = z
}
