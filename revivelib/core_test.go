package revivelib_test

import (
	"strings"
	"testing"

	"github.com/vint-go/vint/config"
	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/revivelib"
	"github.com/vint-go/vint/rules/use_direct_return"
)

func TestReviveLint(t *testing.T) {
	// ARRANGE
	revive := getMockRevive(t)

	// ACT
	failures, err := revive.Lint(revivelib.Include("../rules/use_direct_return/testdata/use_direct_return.go"))
	if err != nil {
		t.Fatal(err)
	}

	// ASSERT
	failureList := []lint.Failure{}

	for failure := range failures {
		failureList = append(failureList, failure)
	}

	const expected = 3

	got := len(failureList)
	if got != expected {
		t.Fatalf("Expected failures to have %d failures, but it has %d.", expected, got)
	}
}

func TestReviveFormat(t *testing.T) {
	t.Setenv("NO_COLOR", "true")
	// ARRANGE
	revive := getMockRevive(t)

	failuresChan, err := revive.Lint(revivelib.Include("../rules/use_direct_return/testdata/use_direct_return.go"))
	if err != nil {
		t.Fatal(err)
	}

	// ACT
	failures, exitCode, err := revive.Format("stylish", failuresChan)
	// ASSERT
	if err != nil {
		t.Fatal(err)
	}

	errorMsgs := []string{
		"(15, 2)  https://revive.run/r#lint/style/useDirectReturn  redundant if ...; err != nil check, just return error instead.",
		"(88, 3)  https://revive.run/r#lint/style/useDirectReturn  redundant if ...; err != nil check, just return error instead.",
		"(95, 3)  https://revive.run/r#lint/style/useDirectReturn  redundant if ...; err != nil check, just return error instead.",
	}
	for _, errorMsg := range errorMsgs {
		if !strings.Contains(failures, errorMsg) {
			t.Fatalf("Expected formatted failures\n'%s'\nto contain\n'%s', but it didn't.", failures, errorMsg)
		}
	}

	const expected = 1
	if exitCode != expected {
		t.Fatalf("Expected exit code to be %d, but it was %d.", expected, exitCode)
	}
}

type mockRule struct{}

func (*mockRule) Name() string {
	return "mock-rule"
}

func (*mockRule) Apply(_ *lint.File, _ lint.Arguments) []lint.Failure {
	return nil
}

func getMockRevive(t *testing.T) *revivelib.Revive {
	t.Helper()

	conf, err := config.GetConfig("../defaults.toml")
	if err != nil {
		t.Fatal(err)
	}

	revive, err := revivelib.New(
		conf,
		true,
		2048,
		revivelib.NewExtraRule(&use_direct_return.IfReturnRule{}, lint.RuleConfig{}),
		revivelib.NewExtraRule(&mockRule{}, lint.RuleConfig{}),
	)
	if err != nil {
		t.Fatal(err)
	}

	return revive
}
