package fixtures

import "fmt"

// Invalid: blank identifier for key with no value variable
func unnecessaryBlankRange(items []string) {
	for _ = range items { // MATCH /unnecessary blank identifier in range statement, use 'for range' instead/
		fmt.Println("item")
	}
}

// Invalid: both key and value are blank identifiers
func bothBlankRange(items []string) {
	for _, _ = range items { // MATCH /unnecessary blank identifiers in range statement, use 'for range' instead/
		fmt.Println("item")
	}
}

// Invalid: blank identifier in channel range
func unnecessaryBlankChannelRange(ch chan int) {
	for _ = range ch { // MATCH /unnecessary blank identifier in range statement, use 'for range' instead/
		fmt.Println("received")
	}
}

// Invalid: blank identifier in map range
func unnecessaryBlankMapRange(m map[string]int) {
	for _ = range m { // MATCH /unnecessary blank identifier in range statement, use 'for range' instead/
		fmt.Println("entry")
	}
}

// Invalid: blank identifier in integer range
func unnecessaryBlankIntRange() {
	for _ = range 10 { // MATCH /unnecessary blank identifier in range statement, use 'for range' instead/
		fmt.Println("count")
	}
}

// Valid: using range without any variables
func validRange(items []string) {
	for range items {
		fmt.Println("item")
	}
}

// Valid: using key variable
func validKeyRange(items []string) {
	for i := range items {
		fmt.Println(i)
	}
}

// Valid: using both key and value variables
func validKeyValueRange(items []string) {
	for i, v := range items {
		fmt.Println(i, v)
	}
}

// Valid: blank key with real value (blank is needed syntactically)
func validBlankKeyWithValue(items []string) {
	for _, v := range items {
		fmt.Println(v)
	}
}

// Valid: blank key with real value in map
func validBlankKeyMapRange(m map[string]int) {
	for _, v := range m {
		fmt.Println(v)
	}
}
