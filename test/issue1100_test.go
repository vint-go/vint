package test_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
)

func TestIssue1100(t *testing.T) {
	testRule(t, "goUnknown/issue1100", &rule.AddConstantRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{
			map[string]any{
				"maxLitCount": "2",
				"allowStrs":   "\"\"",
				"allowInts":   "0,1,2",
				"allowFloats": "0.0,1.0",
				"ignoreFuncs": "os\\.(CreateFile|WriteFile|Chmod|FindProcess),\\.Println,ignoredFunc,\\.Info",
			},
		},
	})
}
