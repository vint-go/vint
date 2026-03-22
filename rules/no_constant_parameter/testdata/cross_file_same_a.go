package fixtures2

import "fmt"

func alwaysFast(mode string) {
	fmt.Println(mode)
}

func callerSameA() {
	alwaysFast("fast")
}
