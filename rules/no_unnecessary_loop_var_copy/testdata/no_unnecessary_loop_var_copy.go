package fixtures

import "fmt"

func rangeLoopSameNameCopy() {
	for i, v := range []int{1, 2, 3} {
		i := i // MATCH /The copy of the 'for' variable "i" can be deleted (Go 1.22+)/
		v := v // MATCH /The copy of the 'for' variable "v" can be deleted (Go 1.22+)/
		fmt.Println(i, v)
	}
}

func cStyleForLoopSameNameCopy() {
	for i := 0; i < 10; i++ {
		i := i // MATCH /The copy of the 'for' variable "i" can be deleted (Go 1.22+)/
		go func() {
			fmt.Println(i)
		}()
	}
}

func cStyleForLoopMultipleInit() {
	for i, j := 1, 1; i+j <= 3; i++ {
		i := i // MATCH /The copy of the 'for' variable "i" can be deleted (Go 1.22+)/
		j := j // MATCH /The copy of the 'for' variable "j" can be deleted (Go 1.22+)/
		fmt.Println(i, j)
	}
}

func multiValueAssignment() {
	for i, v := range []int{1, 2, 3} {
		a, i := 1, i // MATCH /The copy of the 'for' variable "i" can be deleted (Go 1.22+)/
		b, v := 1, v // MATCH /The copy of the 'for' variable "v" can be deleted (Go 1.22+)/
		fmt.Println(a, i, b, v)
	}
}

// Valid: using loop variables directly without copying
func validDirectUse() {
	for i, v := range []int{1, 2, 3} {
		fmt.Println(i, v)
	}
}

// Valid: using loop variable directly in a C-style for loop
func validCStyleDirect() {
	for i := 0; i < 10; i++ {
		go func() {
			fmt.Println(i)
		}()
	}
}

// Valid: creating a new variable with a different value (not a copy of loop var)
func validDifferentValue() {
	for i, v := range []int{1, 2, 3} {
		doubled := v * 2
		idx := i + 1
		fmt.Println(idx, doubled)
	}
}

// Valid: without check-alias, copying to a differently-named variable is allowed by default
func validAliasCopyDefault() {
	for i, v := range []int{1, 2, 3} {
		_i := i
		_v := v
		fmt.Println(_i, _v)
	}
}
