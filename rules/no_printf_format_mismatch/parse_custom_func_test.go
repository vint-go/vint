package no_printf_format_mismatch

import "testing"

func TestParseCustomPrintfFunc_MethodExpression(t *testing.T) {
	pf := parseCustomPrintfFunc("(github.com/my/pkg.Logger).Infof")
	if pf.name != "Infof" {
		t.Errorf("expected name Infof, got %s", pf.name)
	}
	if pf.pkg != "" {
		t.Errorf("expected empty pkg for method expression, got %s", pf.pkg)
	}
	if pf.fmtArgIdx != 0 {
		t.Errorf("expected fmtArgIdx 0, got %d", pf.fmtArgIdx)
	}
	if !pf.isPrintf {
		t.Error("expected isPrintf to be true")
	}
}

func TestParseCustomPrintfFunc_PackageQualified(t *testing.T) {
	pf := parseCustomPrintfFunc("mylog.Debugf")
	if pf.name != "Debugf" {
		t.Errorf("expected name Debugf, got %s", pf.name)
	}
	if pf.pkg != "mylog" {
		t.Errorf("expected pkg mylog, got %s", pf.pkg)
	}
}

func TestParseCustomPrintfFunc_FullyQualifiedPackage(t *testing.T) {
	pf := parseCustomPrintfFunc("github.com/my/pkg/mylog.Debugf")
	if pf.name != "Debugf" {
		t.Errorf("expected name Debugf, got %s", pf.name)
	}
	if pf.pkg != "mylog" {
		t.Errorf("expected pkg mylog (last component), got %s", pf.pkg)
	}
}

func TestParseCustomPrintfFunc_BareFunction(t *testing.T) {
	pf := parseCustomPrintfFunc("myPrintf")
	if pf.name != "myPrintf" {
		t.Errorf("expected name myPrintf, got %s", pf.name)
	}
	if pf.pkg != "" {
		t.Errorf("expected empty pkg, got %s", pf.pkg)
	}
}
