package fixtures

import (
	"reflect"
)

func duplicateArguments() {
	var dst, src []byte

	copy(dst, dst) // MATCH /suspicious duplicated argument: dst appears more than once in the same call/

	reflect.DeepEqual(dst, dst) // MATCH /suspicious duplicated argument: dst appears more than once in the same call/

	// Valid: different arguments
	copy(dst, src)

	// Valid: different arguments
	reflect.DeepEqual(dst, src)

	// Valid: single argument
	_ = len(dst)
}
