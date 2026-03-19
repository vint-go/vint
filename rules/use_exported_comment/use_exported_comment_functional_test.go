package use_exported_comment_test

import (
	"testing"

	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/use_exported_comment"
)

func TestExportedWithDisableStutteringCheck(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_exported_comment_issue_555", &use_exported_comment.ExportedRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"disableStutteringCheck"}})
}

func TestExportedWithChecksOnMethodsOfPrivateTypes(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_exported_comment_issue_552", &use_exported_comment.ExportedRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"checkPrivateReceivers"}})
}

func TestExportedReplacingStuttersByRepetitive(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_exported_comment_issue_519", &use_exported_comment.ExportedRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"sayRepetitiveInsteadOfStutters"}})
}

func TestCheckPublicInterfaceOption(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_exported_comment_issue_1002", &use_exported_comment.ExportedRule{}, &lint.RuleConfig{Arguments: lint.Arguments{"checkPublicInterface"}})
}

func TestCheckDisablingOnDeclarationTypes(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_exported_comment_issue_1045", &use_exported_comment.ExportedRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"disableChecksOnConstants", "disableChecksOnFunctions", "disableChecksOnMethods", "disableChecksOnTypes", "disableChecksOnVariables"},
	})
	functional_test_helpers.TestRule(t, "use_exported_comment_issue_1045", &use_exported_comment.ExportedRule{}, &lint.RuleConfig{
		Arguments: lint.Arguments{"disable-checks-on-constants", "disable-checks-on-functions", "disable-checks-on-methods", "disable-checks-on-types", "disable-checks-on-variables"},
	})
}

func TestCheckDirectiveComment(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_exported_comment_issue_1202", &use_exported_comment.ExportedRule{}, &lint.RuleConfig{})
}

func TestCheckDeprecationComment(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_exported_comment_issue_1231", &use_exported_comment.ExportedRule{}, &lint.RuleConfig{})
}

func TestExportedMainPackage(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_exported_comment_main", &use_exported_comment.ExportedRule{}, &lint.RuleConfig{})
}

func TestCommentVariations(t *testing.T) {
	functional_test_helpers.TestRule(t, "use_exported_comment_issue_1235", &use_exported_comment.ExportedRule{}, &lint.RuleConfig{})
}
