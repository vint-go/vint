package fixtures

import "fmt"

// Invalid: redundant nil check before range over slice
func redundantNilCheckSlice(items []string) {
	if items != nil { // MATCH /redundant nil check on items before range loop; ranging over nil is safe/
		for _, item := range items {
			fmt.Println(item)
		}
	}
}

// Invalid: redundant nil check before range over map
func redundantNilCheckMap(m map[string]int) {
	if m != nil { // MATCH /redundant nil check on m before range loop; ranging over nil is safe/
		for k, v := range m {
			fmt.Println(k, v)
		}
	}
}

// Invalid: redundant nil check with index-only range
func redundantNilCheckIndexRange(xs []int) {
	if xs != nil { // MATCH /redundant nil check on xs before range loop; ranging over nil is safe/
		for range xs {
			fmt.Println("item")
		}
	}
}

// Invalid: redundant nil check with key-only range
func redundantNilCheckKeyRange(xs []int) {
	if xs != nil { // MATCH /redundant nil check on xs before range loop; ranging over nil is safe/
		for i := range xs {
			fmt.Println(i)
		}
	}
}

// Valid: nil check with else branch
func nilCheckWithElse(items []string) {
	if items != nil {
		for _, item := range items {
			fmt.Println(item)
		}
	} else {
		fmt.Println("no items")
	}
}

// Valid: nil check with init statement
func nilCheckWithInit(items []string) {
	if x := items; x != nil {
		for _, item := range x {
			fmt.Println(item)
		}
	}
}

// Valid: nil check with additional statements in body
func nilCheckWithExtraStatements(items []string) {
	if items != nil {
		fmt.Println("processing items")
		for _, item := range items {
			fmt.Println(item)
		}
	}
}

// Valid: nil check on different variable than range target
func nilCheckDifferentVar(items []string, other []string) {
	if items != nil {
		for _, o := range other {
			fmt.Println(o)
		}
	}
}

// Valid: condition is not a nil check
func nonNilCheck(items []string) {
	if len(items) > 0 {
		for _, item := range items {
			fmt.Println(item)
		}
	}
}

// Valid: body has non-range statement only
func bodyHasNonRange(items []string) {
	if items != nil {
		fmt.Println(items)
	}
}

// Valid: nil equality check (== nil), not != nil
func nilEqualCheck(items []string) {
	if items == nil {
		return
	}
	for _, item := range items {
		fmt.Println(item)
	}
}
