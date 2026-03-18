package fixtures

import (
	"bytes"
	"strings"
)

func convenienceFuncExamples(s string, sep string, old string, new string) {
	// Invalid: should use convenience wrappers
	strings.SplitN(s, sep, -1)      // MATCH /strings.SplitN can be replaced with strings.Split/
	strings.SplitAfterN(s, sep, -1) // MATCH /strings.SplitAfterN can be replaced with strings.SplitAfter/
	strings.Replace(s, old, new, -1) // MATCH /strings.Replace can be replaced with strings.ReplaceAll/

	b := []byte(s)
	bsep := []byte(sep)
	bold := []byte(old)
	bnew := []byte(new)

	bytes.SplitN(b, bsep, -1)      // MATCH /bytes.SplitN can be replaced with bytes.Split/
	bytes.SplitAfterN(b, bsep, -1) // MATCH /bytes.SplitAfterN can be replaced with bytes.SplitAfter/
	bytes.Replace(b, bold, bnew, -1) // MATCH /bytes.Replace can be replaced with bytes.ReplaceAll/

	// Valid: these should NOT trigger
	strings.Split(s, sep)
	strings.SplitAfter(s, sep)
	strings.ReplaceAll(s, old, new)

	bytes.Split(b, bsep)
	bytes.SplitAfter(b, bsep)
	bytes.ReplaceAll(b, bold, bnew)

	// Valid: last arg is not -1
	strings.SplitN(s, sep, 2)
	strings.SplitAfterN(s, sep, 0)
	strings.Replace(s, old, new, 1)

	bytes.SplitN(b, bsep, 3)
	bytes.SplitAfterN(b, bsep, 5)
	bytes.Replace(b, bold, bnew, 2)
}
