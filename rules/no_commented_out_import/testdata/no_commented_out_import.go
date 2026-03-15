package fixtures

import (
	"fmt"
	// "os"
	// MATCH /commented-out import: remove or uncomment "os"/
	"strings"
)

import (
	"io"
	// "bytes"
	// MATCH /commented-out import: remove or uncomment "bytes"/
)

// Valid: import block with no commented-out imports
import (
	"path"
	// This is a regular comment describing why we use path
)

// Valid: comments outside import blocks
// "encoding/json"

import "sync" // Valid: single import without parens

func example() {
	fmt.Println("hello")
	_ = strings.NewReader
	_ = io.EOF
	_ = path.Join
	_ = sync.Mutex{}
}
