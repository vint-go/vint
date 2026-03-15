// Test file for noStdlibVersionMismatch rule.
// The testdata/go.mod declares go 1.21, so symbols from go1.22+ should be flagged.

package fixtures

import (
	"cmp"
	"slices"
)

func validUsage() {
	// slices.Sort was introduced in go1.21, which matches the module version.
	s := []int{3, 1, 2}
	slices.Sort(s)
}

func invalidUsage() {
	// slices.Concat was introduced in go1.22, which is newer than go1.21.
	a := []int{1, 2}
	b := []int{3, 4}
	_ = slices.Concat(a, b) // MATCH /slices.Concat requires go1.22 or later (module is go1.21)/

	// cmp.Or was introduced in go1.22, which is newer than go1.21.
	_ = cmp.Or(1, 2) // MATCH /cmp.Or requires go1.22 or later (module is go1.21)/
}
