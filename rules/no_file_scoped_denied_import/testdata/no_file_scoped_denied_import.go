package fixtures

import (
	"fmt"
	"net/http"

	"github.com/stretchr/testify/assert"  // MATCH /import 'github.com/stretchr/testify/assert' is not allowed from list 'main': testify should only be used in test files/
	"github.com/stretchr/testify/require" // MATCH /import 'github.com/stretchr/testify/require' is not allowed from list 'main': testify should only be used in test files/
)

func fileScopedDeniedImportBad() {
	assert.Equal(nil, 1, 1)
	require.Equal(nil, 1, 1)
}

func fileScopedDeniedImportGood() {
	fmt.Println("hello")
	_ = http.StatusOK
}
