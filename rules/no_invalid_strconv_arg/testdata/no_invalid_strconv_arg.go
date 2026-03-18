package fixtures

import "strconv"

func invalidCases() {
	// Invalid base: base must be 0 or between 2 and 36
	_, _ = strconv.ParseInt("42", 1, 64)  // MATCH /invalid base argument to strconv.ParseInt: base must be 0 or between 2 and 36/
	_, _ = strconv.ParseInt("42", 37, 64) // MATCH /invalid base argument to strconv.ParseInt: base must be 0 or between 2 and 36/
	_, _ = strconv.ParseInt("42", -1, 64) // MATCH /invalid base argument to strconv.ParseInt: base must be 0 or between 2 and 36/

	// Invalid bitSize for ParseInt: must be between 0 and 64
	_, _ = strconv.ParseInt("42", 10, 65)  // MATCH /invalid bitSize argument to strconv.ParseInt: bitSize must be between 0 and 64/
	_, _ = strconv.ParseInt("42", 10, -1)  // MATCH /invalid bitSize argument to strconv.ParseInt: bitSize must be between 0 and 64/
	_, _ = strconv.ParseInt("42", 10, 128) // MATCH /invalid bitSize argument to strconv.ParseInt: bitSize must be between 0 and 64/

	// Invalid base for ParseUint
	_, _ = strconv.ParseUint("42", 1, 64)  // MATCH /invalid base argument to strconv.ParseUint: base must be 0 or between 2 and 36/
	_, _ = strconv.ParseUint("42", 37, 64) // MATCH /invalid base argument to strconv.ParseUint: base must be 0 or between 2 and 36/

	// Invalid bitSize for ParseUint
	_, _ = strconv.ParseUint("42", 10, 65) // MATCH /invalid bitSize argument to strconv.ParseUint: bitSize must be between 0 and 64/

	// Invalid bitSize for ParseFloat: must be 32 or 64
	_, _ = strconv.ParseFloat("3.14", 16) // MATCH /invalid bitSize argument to strconv.ParseFloat: bitSize must be 32 or 64/
	_, _ = strconv.ParseFloat("3.14", 0)  // MATCH /invalid bitSize argument to strconv.ParseFloat: bitSize must be 32 or 64/
	_, _ = strconv.ParseFloat("3.14", 96) // MATCH /invalid bitSize argument to strconv.ParseFloat: bitSize must be 32 or 64/

	// Invalid base for FormatInt: must be between 2 and 36
	_ = strconv.FormatInt(42, 0)  // MATCH /invalid base argument to strconv.FormatInt: base must be between 2 and 36/
	_ = strconv.FormatInt(42, 1)  // MATCH /invalid base argument to strconv.FormatInt: base must be between 2 and 36/
	_ = strconv.FormatInt(42, 37) // MATCH /invalid base argument to strconv.FormatInt: base must be between 2 and 36/
}

func validCases() {
	// Valid base 10, bitSize 64
	_, _ = strconv.ParseInt("42", 10, 64)

	// Valid base 0 (auto-detect)
	_, _ = strconv.ParseInt("42", 0, 64)

	// Valid base 2
	_, _ = strconv.ParseInt("101", 2, 32)

	// Valid base 36
	_, _ = strconv.ParseInt("z", 36, 64)

	// Valid bitSize 0 (implies int)
	_, _ = strconv.ParseInt("42", 10, 0)

	// Valid ParseUint
	_, _ = strconv.ParseUint("42", 10, 64)
	_, _ = strconv.ParseUint("42", 0, 32)

	// Valid ParseFloat with bitSize 32
	_, _ = strconv.ParseFloat("3.14", 32)

	// Valid ParseFloat with bitSize 64
	_, _ = strconv.ParseFloat("3.14", 64)

	// Valid FormatInt with various bases
	_ = strconv.FormatInt(42, 2)
	_ = strconv.FormatInt(42, 10)
	_ = strconv.FormatInt(42, 16)
	_ = strconv.FormatInt(42, 36)

	// Variable arguments (not checked, no failure)
	base := 10
	_, _ = strconv.ParseInt("42", base, 64)
}
