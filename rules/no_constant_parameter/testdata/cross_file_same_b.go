package fixtures2

import "fmt"

func callerSameB() {
	alwaysFast("fast") // same constant as file A
	fmt.Println("done")
}
