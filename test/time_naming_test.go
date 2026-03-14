package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestTimeNaming(t *testing.T) {
	testRule(t, "time_naming", &rule.TimeNamingRule{})
}
