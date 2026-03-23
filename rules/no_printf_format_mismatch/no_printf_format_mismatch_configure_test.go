package no_printf_format_mismatch_test

import (
	"testing"

	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/rules/no_printf_format_mismatch"
)

func TestConfigure_WithFuncs(t *testing.T) {
	rule := &no_printf_format_mismatch.NoPrintfFormatMismatchRule{}
	err := rule.Configure(lint.Arguments{
		map[string]any{
			"funcs": []any{
				"(github.com/my/pkg.Logger).Infof",
				"mylog.Debugf",
				"customPrintf",
			},
		},
	})
	if err != nil {
		t.Fatalf("Configure() error: %v", err)
	}
}

func TestConfigure_NoArgs(t *testing.T) {
	rule := &no_printf_format_mismatch.NoPrintfFormatMismatchRule{}
	err := rule.Configure(nil)
	if err != nil {
		t.Fatalf("Configure(nil) error: %v", err)
	}
}

func TestConfigure_EmptyArgs(t *testing.T) {
	rule := &no_printf_format_mismatch.NoPrintfFormatMismatchRule{}
	err := rule.Configure(lint.Arguments{})
	if err != nil {
		t.Fatalf("Configure() error: %v", err)
	}
}

func TestConfigure_NoFuncsKey(t *testing.T) {
	rule := &no_printf_format_mismatch.NoPrintfFormatMismatchRule{}
	err := rule.Configure(lint.Arguments{map[string]any{}})
	if err != nil {
		t.Fatalf("Configure() error: %v", err)
	}
}

func TestConfigure_InvalidArg(t *testing.T) {
	rule := &no_printf_format_mismatch.NoPrintfFormatMismatchRule{}
	err := rule.Configure(lint.Arguments{"not a map"})
	if err == nil {
		t.Fatal("Configure() should return error for invalid argument type")
	}
}

func TestConfigure_InvalidFuncsType(t *testing.T) {
	rule := &no_printf_format_mismatch.NoPrintfFormatMismatchRule{}
	err := rule.Configure(lint.Arguments{map[string]any{
		"funcs": "not a list",
	}})
	if err == nil {
		t.Fatal("Configure() should return error for invalid funcs type")
	}
}

func TestConfigure_InvalidFuncEntry(t *testing.T) {
	rule := &no_printf_format_mismatch.NoPrintfFormatMismatchRule{}
	err := rule.Configure(lint.Arguments{map[string]any{
		"funcs": []any{123},
	}})
	if err == nil {
		t.Fatal("Configure() should return error for invalid func entry type")
	}
}
