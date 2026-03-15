package fixtures

import (
	"fmt"
	printing "fmt" // MATCH /package "fmt" is already imported/
)

import (
	"os"
)

import (
	"strings"
	str "strings" // MATCH /package "strings" is already imported/
)

// Valid: single import of each package is fine.
func singleImport() {
	fmt.Println("hello")
	printing.Println("world")
	_ = os.Stdin
	strings.Contains("a", "b")
	str.Contains("a", "b")
}
