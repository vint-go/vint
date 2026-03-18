package fixtures

import (
	"io/ioutil"
	"io"
)

// Deprecated: Use newHelper instead.
func oldHelper() {}

func newHelper() {}

// Deprecated: Use newConst instead.
const oldConst = 42

const newConst = 42

// Deprecated: Use NewVar instead.
var oldVar = "old"

var newVar = "new"

func example() {
	// Invalid: using a deprecated function from a deprecated package (io/ioutil)
	data, _ := ioutil.ReadAll(nil) // MATCH /io/ioutil.ReadAll is deprecated/
	_ = data

	// Invalid: using a deprecated local function
	oldHelper() // MATCH /fixtures.oldHelper is deprecated/

	// Invalid: using a deprecated local constant
	_ = oldConst // MATCH /fixtures.oldConst is deprecated/

	// Invalid: using a deprecated local variable
	_ = oldVar // MATCH /fixtures.oldVar is deprecated/

	// Valid: using a non-deprecated function from io
	data2, _ := io.ReadAll(nil)
	_ = data2

	// Valid: using a non-deprecated local function
	newHelper()

	// Valid: using non-deprecated constants and variables
	_ = newConst
	_ = newVar
}
