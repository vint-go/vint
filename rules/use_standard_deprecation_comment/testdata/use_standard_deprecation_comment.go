package fixtures

// Invalid examples - should trigger failures

// DEPRECATED: use NewFunc instead
// MATCH /comment looks like a deprecation notice but does not follow the standard "// Deprecated: <explanation>" format/
func OldFunc1() {}

// this function is deprecated, use NewFunc
// MATCH /comment looks like a deprecation notice but does not follow the standard "// Deprecated: <explanation>" format/
func OldFunc2() {}

// deprecated: use NewFunc instead
// MATCH /comment looks like a deprecation notice but does not follow the standard "// Deprecated: <explanation>" format/
func OldFunc3() {}

// Deprecated, use NewFunc instead
// MATCH /comment looks like a deprecation notice but does not follow the standard "// Deprecated: <explanation>" format/
func OldFunc4() {}

// deprecated. use NewFunc instead
// MATCH /comment looks like a deprecation notice but does not follow the standard "// Deprecated: <explanation>" format/
func OldFunc5() {}

// Depricated: use NewFunc instead
// MATCH /comment looks like a deprecation notice but does not follow the standard "// Deprecated: <explanation>" format/
func OldFunc6() {}

// Valid examples - should NOT trigger failures

// Deprecated: use NewFunc instead.
func OldFuncGood() {}

// This is a normal comment, nothing about deprecation.
func NormalFunc() {}

func NoComment() {}
