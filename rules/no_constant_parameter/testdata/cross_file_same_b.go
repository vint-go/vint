package fixtures2

import "fmt"

func callerSameB() {
	alwaysFast("fast") // same constant as file A
	alwaysFast("fast")
	fmt.Println("done")
}
