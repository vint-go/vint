package example

import "os"

func doSomething() {
	os.Open("file.txt") //nolint:errcheck
}
