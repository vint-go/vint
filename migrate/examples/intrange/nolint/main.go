package example

import "fmt"

func loopExample() {
	for i := 0; i < 10; i++ { //nolint:intrange
		fmt.Println(i)
	}
}
