package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestJsonDataFormatVarNaming(t *testing.T) {
	testRule(t, "json_data_format_var_naming", &rule.VarNamingRule{}, &lint.RuleConfig{})
}
