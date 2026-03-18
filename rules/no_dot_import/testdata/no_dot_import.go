package fixtures

import . "fmt" // MATCH /should not use dot imports/

import "os"

import (
	. "strings" // MATCH /should not use dot imports/
	"io"
)

// Valid: regular imports are fine.
func good() {
	_ = os.Stdout
	_ = io.EOF
}
