package fixtures

import (
	"bytes"
	"strings"
)

func swappedArgs() {
	s := "hello world"
	b := []byte("hello world")

	// Invalid: arguments appear to be swapped
	strings.HasPrefix("prefix", s) // MATCH /strings.HasPrefix arguments order looks suspicious, the literal string argument should probably not be first/
	strings.HasSuffix("suffix", s) // MATCH /strings.HasSuffix arguments order looks suspicious, the literal string argument should probably not be first/
	strings.Contains("needle", s)  // MATCH /strings.Contains arguments order looks suspicious, the literal string argument should probably not be first/
	strings.EqualFold("foo", s)    // MATCH /strings.EqualFold arguments order looks suspicious, the literal string argument should probably not be first/
	strings.Split("sep", s)        // MATCH /strings.Split arguments order looks suspicious, the literal string argument should probably not be first/
	strings.SplitAfter("sep", s)   // MATCH /strings.SplitAfter arguments order looks suspicious, the literal string argument should probably not be first/
	strings.SplitN("sep", s, 2)    // MATCH /strings.SplitN arguments order looks suspicious, the literal string argument should probably not be first/
	strings.SplitAfterN("sep", s, 2) // MATCH /strings.SplitAfterN arguments order looks suspicious, the literal string argument should probably not be first/
	strings.Count("sub", s)        // MATCH /strings.Count arguments order looks suspicious, the literal string argument should probably not be first/
	strings.Index("sub", s)        // MATCH /strings.Index arguments order looks suspicious, the literal string argument should probably not be first/
	strings.Replace("old", s, "new", -1) // MATCH /strings.Replace arguments order looks suspicious, the literal string argument should probably not be first/
	strings.ReplaceAll("old", s, "new")  // MATCH /strings.ReplaceAll arguments order looks suspicious, the literal string argument should probably not be first/
	strings.TrimPrefix("prefix", s) // MATCH /strings.TrimPrefix arguments order looks suspicious, the literal string argument should probably not be first/
	strings.TrimSuffix("suffix", s) // MATCH /strings.TrimSuffix arguments order looks suspicious, the literal string argument should probably not be first/
	strings.Cut("sep", s)          // MATCH /strings.Cut arguments order looks suspicious, the literal string argument should probably not be first/

	bytes.HasPrefix([]byte("prefix"), b) // bytes uses []byte, not string literal, so no match
	bytes.Contains([]byte("needle"), b)  // bytes uses []byte, not string literal, so no match

	// Valid: correct argument order
	strings.HasPrefix(s, "prefix")
	strings.HasSuffix(s, "suffix")
	strings.Contains(s, "needle")
	strings.EqualFold(s, "foo")
	strings.Split(s, "sep")
	strings.SplitAfter(s, "sep")
	strings.SplitN(s, "sep", 2)
	strings.SplitAfterN(s, "sep", 2)
	strings.Count(s, "sub")
	strings.Index(s, "sub")
	strings.Replace(s, "old", "new", -1)
	strings.ReplaceAll(s, "old", "new")
	strings.TrimPrefix(s, "prefix")
	strings.TrimSuffix(s, "suffix")
	strings.Cut(s, "sep")
	bytes.HasPrefix(b, []byte("prefix"))
	bytes.Contains(b, []byte("needle"))

	// Valid: both arguments are literals
	strings.HasPrefix("hello", "he")
	strings.Contains("hello", "ell")

	// Valid: both arguments are variables
	s2 := "world"
	strings.HasPrefix(s, s2)

	_ = s
	_ = b
}
