//go:build linux

//go:build linux
// MATCH /duplicate build constraint: "//go:build linux"/

package fixtures

// Valid: single build constraint is fine.

func example() {}
