package fixtures

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/sirupsen/logrus"  // MATCH /import 'github.com/sirupsen/logrus' is not allowed from list 'main': Use log/slog instead of logrus/
	"github.com/pkg/errors"       // MATCH /import 'github.com/pkg/errors' is not allowed from list 'main': Use fmt.Errorf with %w verb for error wrapping/
	"io/ioutil"                   // MATCH /import 'io/ioutil' is not allowed from list 'main': Use os and io packages directly instead/
)

func deniedImportBad() {
	logrus.Info("hello")
	err := errors.New("something failed")
	_ = err
	data, _ := ioutil.ReadFile("test.txt")
	_ = data
}

func deniedImportGood() {
	slog.Info("hello")
	err := fmt.Errorf("something failed: %w", fmt.Errorf("cause"))
	_ = err
	data, _ := os.ReadFile("test.txt")
	_ = data
}
