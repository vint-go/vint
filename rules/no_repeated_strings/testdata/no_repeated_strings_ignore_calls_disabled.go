package fixtures

import "fmt"

// Invalid: With ignore-calls disabled, strings appearing only in
// function call arguments are also flagged.
func callsOnlyFlagged() {
	fmt.Println("call_only_str") // MATCH /string literal "call_only_str" appears 3 times, consider extracting it into a named constant/
	fmt.Println("call_only_str")
	fmt.Println("call_only_str")
}

// Invalid: Strings in non-call contexts are still flagged.
func nonCallContext() {
	x := "assigned_str" // MATCH /string literal "assigned_str" appears 3 times, consider extracting it into a named constant/
	y := "assigned_str"
	_ = x == "assigned_str"
	_ = y
}
