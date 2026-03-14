package fixtures

import "fmt"

var (
	nutaFoo  any = "foo"
	nutaBars     = []any{1, 42, "some", "thing"}
)

// Invalid: type assertion without checking the ok value
func nutaAssignWithoutOk() {
	s := nutaFoo.(string) // MATCH /type assertion result is unchecked in nutaFoo.(string), use the comma-ok idiom/
	fmt.Println(s)
}

// Invalid: type assertion in function argument without checking
func nutaInFunctionArg() {
	fmt.Println(nutaFoo.(int) + 1) // MATCH /type assertion result is unchecked in nutaFoo.(int), use the comma-ok idiom/
}

// Invalid: type assertion in return
func nutaReturn() int {
	return nutaFoo.(int) // MATCH /type assertion result is unchecked in nutaFoo.(int), use the comma-ok idiom/
}

// Invalid: type assertion in switch tag
func nutaSwitch() {
	switch nutaFoo.(int) { // MATCH /type assertion result is unchecked in nutaFoo.(int), use the comma-ok idiom/
	case 0:
	case 1:
	}
}

// Invalid: type assertion in range
func nutaRange() {
	var some any = nutaBars
	for _, x := range some.([]string) { // MATCH /type assertion result is unchecked in some.([]string), use the comma-ok idiom/
		fmt.Println(x)
	}
}

// Invalid: type assertion in if condition
func nutaIfCondition() {
	if nutaFoo.(int) == 1 { // MATCH /type assertion result is unchecked in nutaFoo.(int), use the comma-ok idiom/
		return
	}
}

// Invalid: type assertion in if condition (reverse)
func nutaIfConditionReverse() {
	if 1 == nutaFoo.(int) { // MATCH /type assertion result is unchecked in nutaFoo.(int), use the comma-ok idiom/
		return
	}
}

// Invalid: type assertion with ok value ignored via _
func nutaIgnoredOk() {
	r, _ := nutaFoo.(int) // MATCH /type assertion result is unchecked in nutaFoo.(int), use the comma-ok idiom/
	fmt.Println(r)
}

// Invalid: type assertion in channel send
func nutaChannelSend() {
	c := make(chan any)
	var a any = "foo"
	c <- a.(int) // MATCH /type assertion result is unchecked in a.(int), use the comma-ok idiom/
}

// Invalid: type assertion in switch comparison
func nutaSwitchComparison() {
	switch nutaFoo.(int) == 1 { // MATCH /type assertion result is unchecked in nutaFoo.(int), use the comma-ok idiom/
	case true:
	case false:
	}
}

// Invalid: type assertion in case clause
func nutaCaseClause() {
	switch {
	case nutaFoo.(int) == 1: // MATCH /type assertion result is unchecked in nutaFoo.(int), use the comma-ok idiom/
	}
}

// Valid: type assertion using the comma-ok idiom
func nutaCommaOk() {
	s, ok := nutaFoo.(string)
	if !ok {
		fmt.Println("not a string")
		return
	}
	fmt.Println(s)
}

// Valid: type assertion with error handling
func nutaCommaOkProcess() {
	num, ok := nutaFoo.(int)
	if !ok {
		fmt.Println("expected an integer")
		return
	}
	fmt.Println(num + 1)
}

// Valid: type switch (not a direct type assertion)
func nutaTypeSwitch() {
	switch v := nutaFoo.(type) {
	case string:
		fmt.Printf("String: %s\n", v)
	case int:
		fmt.Printf("Integer: %d\n", v)
	default:
		fmt.Printf("Unknown type: %T\n", v)
	}
}

// Valid: type switch without assignment
func nutaTypeSwitchNoAssign() {
	switch nutaFoo.(type) {
	case string:
	case int:
	}
}
