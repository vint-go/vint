package fixtures

import (
	"fmt"
	"net/http"

	"github.com/myorg/service"
	"github.com/some/other-lib"    // MATCH /import 'github.com/some/other-lib' is not allowed from list 'main'/
	"github.com/thirdparty/denied" // MATCH /import 'github.com/thirdparty/denied' is not allowed from list 'main'/
)

func unallowedImportBad() {
	other.Use()
	denied.Use()
}

func unallowedImportGood() {
	fmt.Println("hello")
	_ = http.StatusOK
	service.Start()
}
