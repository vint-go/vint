package fixtures

import (
	"errors"
	"fmt"
	"sort"
)

func unusedFmtSprintf() {
	fmt.Sprintf("hello %s", "world") // MATCH /result of fmt.Sprintf call not used/
}

func unusedFmtSprint() {
	fmt.Sprint("hello") // MATCH /result of fmt.Sprint call not used/
}

func unusedFmtErrorf() {
	fmt.Errorf("something %s", "bad") // MATCH /result of fmt.Errorf call not used/
}

func unusedErrorsNew() {
	errors.New("something went wrong") // MATCH /result of errors.New call not used/
}

func unusedSortReverse() {
	s := []int{3, 1, 2}
	sort.Reverse(sort.IntSlice(s)) // MATCH /result of sort.Reverse call not used/
}

func usedFmtSprintf() string {
	// Good: result is assigned to a variable
	msg := fmt.Sprintf("hello %s", "world")
	return msg
}

func usedFmtSprintReturn() string {
	// Good: result is returned directly
	return fmt.Sprintf("hello %s", "world")
}

func usedErrorsNew() error {
	// Good: result is returned
	return errors.New("something went wrong")
}

func usedFmtErrorf() error {
	// Good: result is assigned
	err := fmt.Errorf("bad %s", "thing")
	return err
}

func usedSortReverse() {
	s := []int{3, 1, 2}
	// Good: result is used
	r := sort.Reverse(sort.IntSlice(s))
	sort.Sort(r)
}

func fmtPrintlnOk() {
	// Good: fmt.Println is not a pure function (has side effects)
	fmt.Println("hello")
}
