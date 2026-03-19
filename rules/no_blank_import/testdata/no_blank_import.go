// Package fixtures tests that blank imports with justifying comments are allowed.
package fixtures

// OK

import (
	"fmt"

	// This blank import is justified with a doc comment.
	_ "os"

	_ "strings" // This blank import is justified with an inline comment.
)

func _() {
	fmt.Println("use fmt")
}
