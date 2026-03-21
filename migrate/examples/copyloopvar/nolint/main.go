package example

import "fmt"

func process() {
	for i, v := range []int{1, 2, 3} {
		i := i //nolint:copyloopvar
		v := v //nolint:copyloopvar
		fmt.Println(i, v)
	}
}
