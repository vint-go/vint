package fixtures

import "fmt"

func badSprintf(dynamicString string) {
	_ = fmt.Sprintf(dynamicString) // MATCH /dynamic string used as format argument without additional args; use fmt.Sprint or add format arguments/
}

func badPrintf(dynamicString string) {
	fmt.Printf(dynamicString) // MATCH /dynamic string used as format argument without additional args; use fmt.Print or add format arguments/
}

func badErrorf(dynamicString string) {
	_ = fmt.Errorf(dynamicString) // MATCH /dynamic string used as format argument without additional args; use fmt.Error or add format arguments/
}

func goodDirect(dynamicString string) {
	_ = dynamicString
}

func goodWithFormatVerb(dynamicString string) {
	_ = fmt.Sprintf("%s", dynamicString)
}

func goodWithPrefix(dynamicString string) {
	_ = fmt.Sprintf("prefix: %s", dynamicString)
}

func goodLiteralFormat() {
	_ = fmt.Sprintf("hello world")
}

func goodSprintfMultipleArgs(dynamicString string) {
	_ = fmt.Sprintf(dynamicString, 42)
}

func goodPrintfMultipleArgs(dynamicString string) {
	fmt.Printf(dynamicString, "value")
}
