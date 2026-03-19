package test_test

import (
	"flag"
	"os"
	"path/filepath"
	"regexp"
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
	"github.com/strowk/vint/rules/no_blank_import"
	"github.com/strowk/vint/rules/use_var_naming"
	"github.com/strowk/vint/rules/no_var_declaration"
	"github.com/strowk/vint/rules/no_context_keys_type"
	"github.com/strowk/vint/rules/no_dot_import"
	"github.com/strowk/vint/rules/use_error_last_return"
	"github.com/strowk/vint/rules/no_error_strings"
	"github.com/strowk/vint/rules/use_errorf"
	"github.com/strowk/vint/rules/use_error_naming"
	"github.com/strowk/vint/rules/use_exported_comment"
	"github.com/strowk/vint/rules/use_indent_error_flow"
	"github.com/strowk/vint/rules/use_package_comments"
	"github.com/strowk/vint/rules/no_redundant_range_val"
	"github.com/strowk/vint/rules/use_receiver_naming"
	"github.com/strowk/vint/rules/no_unexported_return"
	"github.com/strowk/vint/rules/use_time_naming"
)

var lintMatch = flag.String("lint.match", "", "restrict fixtures matches to this pattern")

var rules = []lint.Rule{
	&no_var_declaration.VarDeclarationsRule{},
	&use_package_comments.PackageCommentsRule{},
	&no_dot_import.NoDotImportRule{},
	&no_blank_import.NoBlankImportRule{},
	&use_exported_comment.ExportedRule{},
	&use_var_naming.VarNamingRule{},
	&use_indent_error_flow.IndentErrorFlowRule{},
	&no_redundant_range_val.RangeRule{},
	&use_errorf.ErrorfRule{},
	&use_error_naming.ErrorNamingRule{},
	&no_error_strings.ErrorStringsRule{},
	&use_receiver_naming.ReceiverNamingRule{},
	&rule.IncrementDecrementRule{},
	&use_error_last_return.UseErrorLastReturnRule{},
	&no_unexported_return.UnexportedReturnRule{},
	&use_time_naming.TimeNamingRule{},
	&no_context_keys_type.ContextKeysType{},
}

func TestAll(t *testing.T) {
	baseDir := "../testdata/golint/"

	for _, r := range rules {
		configureRule(t, r, nil)
	}

	rx, err := regexp.Compile(*lintMatch)
	if err != nil {
		t.Fatalf("Bad -lint.match value %q: %v", *lintMatch, err)
	}

	fis, err := os.ReadDir(baseDir)
	if err != nil {
		t.Fatalf("os.ReadDir: %v", err)
	}
	if len(fis) == 0 {
		t.Fatalf("no files in %v", baseDir)
	}
	for _, fi := range fis {
		if !rx.MatchString(fi.Name()) {
			continue
		}
		t.Run(fi.Name(), func(t *testing.T) {
			filePath := filepath.Join(baseDir, fi.Name())
			src, err := os.ReadFile(filePath) //nolint:gosec // ignore G304: potential file inclusion via variable
			if err != nil {
				t.Fatalf("Failed reading %s: %v", fi.Name(), err)
			}

			ins := parseInstructions(t, filePath, src)

			assertFailures(t, filePath, rules, map[string]lint.RuleConfig{}, ins)
		})
	}
}
