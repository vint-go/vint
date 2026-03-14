package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestJsonDataFormat(t *testing.T) {
	testRule(t, "json_data_format_atomic", &rule.AtomicRule{})
}

func TestJsonDataFormatVarNaming(t *testing.T) {
	testRule(t, "json_data_format_var_naming", &rule.VarNamingRule{}, &lint.RuleConfig{})
}
