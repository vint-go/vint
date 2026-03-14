package no_insecure_random

import (
	"go/ast"
	"strconv"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoInsecureRandomRule detects the use of insecure random number sources from the math/rand package.
type NoInsecureRandomRule struct{}

// Apply applies the rule to given file.
func (r *NoInsecureRandomRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Find the local name used for math/rand import
	randPkgName := mathRandImportName(file.AST)
	if randPkgName == "" {
		return failures
	}

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintInsecureRandom{
		onFailure:   onFailure,
		randPkgName: randPkgName,
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoInsecureRandomRule) Name() string {
	return "noInsecureRandom"
}

// Group returns the rule group.
func (*NoInsecureRandomRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoInsecureRandomRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// mathRandImportName returns the local name used for the "math/rand" import,
// or empty string if math/rand is not imported.
func mathRandImportName(f *ast.File) string {
	for _, imp := range f.Imports {
		path, err := strconv.Unquote(imp.Path.Value)
		if err != nil {
			continue
		}
		if path == "math/rand" || path == "math/rand/v2" {
			if imp.Name != nil {
				if imp.Name.Name == "_" || imp.Name.Name == "." {
					return ""
				}
				return imp.Name.Name
			}
			return "rand"
		}
	}
	return ""
}

// insecureRandomFuncs is the set of math/rand functions that produce random values
// and are considered insecure for security-sensitive purposes.
var insecureRandomFuncs = map[string]bool{
	"Int":         true,
	"Intn":        true,
	"Int31":       true,
	"Int31n":      true,
	"Int63":       true,
	"Int63n":      true,
	"Int32":       true,
	"Int32N":      true,
	"Int64":       true,
	"Int64N":      true,
	"IntN":        true,
	"Uint":        true,
	"Uint32":      true,
	"Uint32N":     true,
	"Uint64":      true,
	"Uint64N":     true,
	"UintN":       true,
	"Float32":     true,
	"Float64":     true,
	"NormFloat64": true,
	"ExpFloat64":  true,
	"Perm":        true,
	"Read":        true,
	"N":           true,
}

type lintInsecureRandom struct {
	onFailure   func(lint.Failure)
	randPkgName string
}

func (w *lintInsecureRandom) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	sel, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return w
	}

	if ident.Name != w.randPkgName {
		return w
	}

	funcName := sel.Sel.Name
	if !insecureRandomFuncs[funcName] {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       ce,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "use of insecure random number generator (math/rand), use crypto/rand instead",
	})

	return w
}
