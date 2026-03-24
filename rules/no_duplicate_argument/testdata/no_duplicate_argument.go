package fixtures

import (
	"bytes"
	"reflect"
	"strings"
)

func duplicateArguments() {
	var dst, src []byte
	var s, t string

	copy(dst, dst) // MATCH /suspicious duplicated argument: dst appears more than once in the same call/

	reflect.DeepEqual(dst, dst) // MATCH /suspicious duplicated argument: dst appears more than once in the same call/

	strings.Contains(s, s) // MATCH /suspicious duplicated argument: s appears more than once in the same call/

	strings.Replace(s, t, t, -1) // MATCH /suspicious duplicated argument: t appears more than once in the same call/

	bytes.Equal(dst, dst) // MATCH /suspicious duplicated argument: dst appears more than once in the same call/

	// Valid: different arguments
	copy(dst, src)

	// Valid: different arguments
	reflect.DeepEqual(dst, src)

	// Valid: single argument
	_ = len(dst)

	// Valid: not a whitelisted function (arbitrary function calls should not be checked)
	someFunc(s, s)

	// Valid: not a whitelisted function
	sqlmock.NewResult(0, 0)

	// Valid: impure expressions (function calls)
	reflect.DeepEqual(newThing(), newThing())
}

func someFunc(a, b string) {}

func newThing() interface{} { return nil }

type sqlmockPkg struct{}

var sqlmock sqlmockPkg

func (sqlmockPkg) NewResult(a, b int64) {}
