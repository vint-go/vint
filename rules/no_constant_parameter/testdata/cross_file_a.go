package fixtures

import "fmt"

// process is defined here and called from both files.
func process(mode string) {
	fmt.Println(mode)
}

func callerA() {
	process("fast") // file A always passes "fast"
}
