package no_stdlib_version_mismatch

import (
	"fmt"
	"go/ast"
	"go/types"
	"go/version"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoStdlibVersionMismatchRule reports uses of standard library symbols that are
// "too new" for the Go version in effect. If your go.mod file specifies go 1.21
// but your code uses a function introduced in Go 1.22, this rule will flag it.
type NoStdlibVersionMismatchRule struct{}

// Apply applies the rule to the given file.
func (r *NoStdlibVersionMismatchRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	goVer := file.Pkg.GoVersionString()
	if goVer == "" {
		return nil
	}

	// Skip modules before go1.21 since the go directive wasn't clearly
	// specified as a minimum version requirement before that release.
	if version.Compare(version.Lang(goVer), "go1.21") < 0 {
		return nil
	}

	if err := file.Pkg.TypeCheck(); err != nil {
		// Type checking may produce errors but still provides partial info.
		_ = err
	}

	typesInfo := file.Pkg.TypesInfo()
	if typesInfo == nil {
		return nil
	}

	// Compute the module's minor Go version number for fast comparison.
	// goVer is in "go1.X" format.
	moduleMinor := goMinorVersion(goVer)
	if moduleMinor < 21 {
		return nil
	}

	w := &lintStdlibVersion{
		moduleGoVer: goVer,
		moduleMinor: moduleMinor,
		typesInfo:   typesInfo,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}

	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoStdlibVersionMismatchRule) Name() string {
	return "noStdlibVersionMismatch"
}

// Group returns the rule group.
func (*NoStdlibVersionMismatchRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoStdlibVersionMismatchRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func (*NoStdlibVersionMismatchRule) RequiresTypecheck() bool {
	return true
}

type lintStdlibVersion struct {
	moduleGoVer string
	moduleMinor int
	typesInfo   *types.Info
	onFailure   func(lint.Failure)
}

func (w *lintStdlibVersion) Visit(node ast.Node) ast.Visitor {
	ident, ok := node.(*ast.Ident)
	if !ok {
		return w
	}

	obj, ok := w.typesInfo.Uses[ident]
	if !ok || obj.Pkg() == nil {
		return w
	}

	pkgPath := obj.Pkg().Path()
	if !stdlibPackages[pkgPath] {
		return w
	}

	// Check if this symbol is too new for the module's Go version.
	symKey := pkgPath + "." + obj.Name()
	minMinor, found := stdlibSymbolVersions[symKey]
	if !found {
		return w
	}

	if w.moduleMinor < minMinor {
		minVersion := fmt.Sprintf("go1.%d", minMinor)
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       ident,
			Failure: fmt.Sprintf("%s.%s requires %s or later (module is %s)",
				obj.Pkg().Name(), obj.Name(), minVersion, w.moduleGoVer),
		})
	}

	return w
}

// goMinorVersion extracts the minor version number from a "go1.X" string.
// Returns 0 if the format is invalid.
func goMinorVersion(goVer string) int {
	// goVer is in "go1.X" format, e.g. "go1.21"
	var major, minor int
	n, _ := fmt.Sscanf(goVer, "go%d.%d", &major, &minor)
	if n == 2 && major == 1 {
		return minor
	}
	return 0
}
