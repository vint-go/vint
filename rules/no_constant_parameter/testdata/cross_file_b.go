package fixtures

import "fmt"

func callerB() {
	process("slow") // file B passes "slow" — different from file A
	fmt.Println("done")
}
