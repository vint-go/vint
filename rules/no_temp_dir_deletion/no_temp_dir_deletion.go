package no_temp_dir_deletion

import (
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoTempDirDeletionRule detects calls that delete the system temp directory
// returned by os.TempDir(), which can affect other processes.
type NoTempDirDeletionRule struct{}

// Apply applies the rule to given file.
func (r *NoTempDirDeletionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintTempDirDeletion{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoTempDirDeletionRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintTempDirDeletion{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoTempDirDeletionRule) Name() string {
	return "noTempDirDeletion"
}

// Group returns the rule group.
func (*NoTempDirDeletionRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoTempDirDeletionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// removeFunctions lists the os functions that delete files/directories.
var removeFunctions = []struct {
	pkg  string
	name string
}{
	{"os", "Remove"},
	{"os", "RemoveAll"},
}

type lintTempDirDeletion struct {
	onFailure func(lint.Failure)
}

func (w *lintTempDirDeletion) Visit(node ast.Node) ast.Visitor {
	ce, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	for _, fn := range removeFunctions {
		if !astutils.IsPkgDotName(ce.Fun, fn.pkg, fn.name) {
			continue
		}
		if len(ce.Args) < 1 {
			continue
		}
		if isOsTempDirCall(ce.Args[0]) {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       ce,
				Category:   lint.FailureCategoryLogic,
				Failure:    "deleting the system temp directory affects other processes; use a subdirectory instead",
			})
		}
	}

	return w
}

// isOsTempDirCall checks whether the given expression is a call to os.TempDir().
func isOsTempDirCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	return astutils.IsPkgDotName(call.Fun, "os", "TempDir")
}
