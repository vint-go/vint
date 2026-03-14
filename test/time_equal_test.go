package test_test

import (
	"testing"

	"github.com/strowk/vint/rule"
)

func TestTimeEqual(t *testing.T) {
	testRule(t, "time_equal", &rule.TimeEqualRule{})
}
