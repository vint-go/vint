//go:build linux AND darwin
// MATCH:1 /malformed build tag: unexpected token AND/

package fixtures

//go:build linux
// MATCH:6 /build tag must appear before the package clause/

// Valid examples below - no matches expected

func example() {}
