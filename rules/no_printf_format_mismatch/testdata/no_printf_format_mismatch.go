package fixtures

import (
	"fmt"
	"log"
)

// Invalid: %d expects an integer, got a string
func badVerbType() {
	name := "world"
	fmt.Printf("Hello %d", name) // MATCH /format verb %d expects an integer type, got string (argument #1)/
}

// Invalid: too few arguments for format string
func tooFewArgs() {
	fmt.Printf("Name: %s, Age: %d", "Alice") // MATCH /format string expects 2 argument(s), but 1 provided/
}

// Invalid: Println does not support format verbs
func printlnWithVerb() {
	fmt.Println("Hello %s", "world") // MATCH /fmt.Println call has possible formatting directive; use fmt.Printlnf instead/
}

// Invalid: too many arguments
func tooManyArgs() {
	fmt.Printf("Hello %s", "world", "extra") // MATCH /format string expects 1 argument(s), but 2 provided/
}

// Invalid: Sprintf with wrong count
func sprintfWrongCount() {
	_ = fmt.Sprintf("a=%d b=%d", 1) // MATCH /format string expects 2 argument(s), but 1 provided/
}

// Invalid: log.Printf mismatch
func logPrintfMismatch() {
	log.Printf("value: %d %s", 42) // MATCH /format string expects 2 argument(s), but 1 provided/
}

// Invalid: Print with format verb
func printWithVerb() {
	fmt.Print("value: %d", 42) // MATCH /fmt.Print call has possible formatting directive; use fmt.Printf instead/
}

// Valid: correct format usage
func validPrintf() {
	name := "world"
	fmt.Printf("Hello %s", name)
}

// Valid: correct format with multiple args
func validMultiArg() {
	fmt.Printf("Name: %s, Age: %d", "Alice", 30)
}

// Valid: Println without format verbs
func validPrintln() {
	fmt.Println("Hello", "world")
}

// Valid: no arguments to Println
func validPrintlnNoArgs() {
	fmt.Println("Hello world")
}

// Valid: Sprintf with correct args
func validSprintf() {
	_ = fmt.Sprintf("a=%d b=%s", 1, "two")
}

// Valid: percent literal
func validPercentLiteral() {
	fmt.Printf("100%% done")
}

// Valid: log.Printf correct
func validLogPrintf() {
	log.Printf("value: %d", 42)
}

// Valid: %v accepts any type
func validVerbV() {
	fmt.Printf("value: %v", "anything")
	fmt.Printf("value: %v", 42)
	fmt.Printf("value: %v", true)
}
