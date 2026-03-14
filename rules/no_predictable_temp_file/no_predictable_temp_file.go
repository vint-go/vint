package no_predictable_temp_file

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoPredictableTempFileRule detects the creation of temporary files using
// predictable paths instead of using os.CreateTemp or ioutil.TempFile.
type NoPredictableTempFileRule struct{}

// Apply applies the rule to given file.
func (r *NoPredictableTempFileRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintPredictableTempFile{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoPredictableTempFileRule) Name() string {
	return "noPredictableTempFile"
}

// Group returns the rule group.
func (*NoPredictableTempFileRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoPredictableTempFileRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// tempDirPrefixes lists the common temporary directory prefixes to check.
var tempDirPrefixes = []string{
	"/tmp/",
	"/var/tmp/",
	"/dev/shm/",
}

// fileCreationFuncs maps package.Function to the index of the path argument.
var fileCreationFuncs = map[string]map[string]int{
	"os": {
		"Create":    0,
		"OpenFile":  0,
		"WriteFile": 0,
		"Mkdir":     0,
		"MkdirAll":  0,
	},
	"ioutil": {
		"WriteFile": 0,
	},
}

type lintPredictableTempFile struct {
	onFailure func(lint.Failure)
}

func (w *lintPredictableTempFile) Visit(node ast.Node) ast.Visitor {
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

	pkgFuncs, ok := fileCreationFuncs[ident.Name]
	if !ok {
		return w
	}

	argIdx, ok := pkgFuncs[sel.Sel.Name]
	if !ok {
		return w
	}

	if len(ce.Args) <= argIdx {
		return w
	}

	pathArg := ce.Args[argIdx]
	lit, ok := pathArg.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return w
	}

	pathVal, err := strconv.Unquote(lit.Value)
	if err != nil {
		return w
	}

	if hasTempDirPrefix(pathVal) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "predictable temporary file path, use os.CreateTemp instead",
		})
	}

	return w
}

// hasTempDirPrefix checks if the given path starts with a known temporary directory prefix.
func hasTempDirPrefix(path string) bool {
	lower := strings.ToLower(path)
	for _, prefix := range tempDirPrefixes {
		if strings.HasPrefix(lower, prefix) {
			return true
		}
	}
	return false
}
