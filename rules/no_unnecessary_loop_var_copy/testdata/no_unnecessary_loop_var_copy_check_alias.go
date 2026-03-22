package fixtures

import "fmt"

// With check-alias enabled: copying loop variable to a differently-named variable
func aliasRangeLoop() {
	for i, v := range []int{1, 2, 3} {
		_i := i // MATCH /The copy of the 'for' variable "i" can be deleted (Go 1.22+)/
		_v := v // MATCH /The copy of the 'for' variable "v" can be deleted (Go 1.22+)/
		fmt.Println(_i, _v)
	}
}

// With check-alias enabled: multi-value assignment with aliased loop variable
func aliasMultiValue() {
	for i := range []int{1, 2, 3} {
		b, _i := 1, i // MATCH /The copy of the 'for' variable "i" can be deleted (Go 1.22+)/
		fmt.Println(b, _i)
	}
}

// Same-name copies should still be flagged with check-alias enabled
func sameNameWithCheckAlias() {
	for i, v := range []int{1, 2, 3} {
		i := i // MATCH /The copy of the 'for' variable "i" can be deleted (Go 1.22+)/
		v := v // MATCH /The copy of the 'for' variable "v" can be deleted (Go 1.22+)/
		fmt.Println(i, v)
	}
}

// Valid: using loop variables directly without copying
func validDirectUseAlias() {
	for i, v := range []int{1, 2, 3} {
		fmt.Println(i, v)
	}
}

// Valid: creating a new variable with a different value (not a copy of loop var)
func validDifferentValueAlias() {
	for i, v := range []int{1, 2, 3} {
		doubled := v * 2
		idx := i + 1
		fmt.Println(idx, doubled)
	}
}
