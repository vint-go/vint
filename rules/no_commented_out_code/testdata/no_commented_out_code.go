package fixtures

import "fmt"

func process() {
	// fmt.Println("debug output")
	// MATCH /comment appears to contain commented-out code: fmt.Println("debug output")/

	doWork()
}

func assignFunc() {
	// result := computeValue()
	// MATCH /comment appears to contain commented-out code: result := computeValue()/

	doWork()
}

func controlFlowFunc() {
	// if err != nil { return err }
	// MATCH /comment appears to contain commented-out code: if err != nil { return err }/

	doWork()
}

// Valid examples - should NOT trigger failures

func goodFunc() {
	// Process the work items and return results
	doWork()
}

func goodFunc2() {
	// TODO: refactor this later
	doWork()
}

func goodFunc3() {
	// See https://example.com for details
	doWork()
}

func goodFunc4() {
	// e.g. SomeType is used for something
	doWork()
}

func goodFunc5() {
	// short
	doWork()
}

func goodFunc6() {
	// SomeType is a reference to documentation
	doWork()
}

func goodFunc7() {
	// TODO: 404
	// fmt.Println("debug output")
	// result := computeValue()
	doWork()
}

func goodFunc8() {
	// See https://example.com/issue/123
	// if err != nil { return err }
	doWork()
}

// Top-level commented-out code should NOT trigger (only inside functions)
// fmt.Println("top level debug")

func doWork()       {}
func computeValue() int { return 0 }

func unused() {
	_ = fmt.Sprintf("")
}
