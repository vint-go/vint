package no_import_shadow

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoImportShadowRule detects when imported package names are shadowed in assignments.
type NoImportShadowRule struct{}

// Apply applies the rule to given file.
func (r *NoImportShadowRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	importNames := map[string]struct{}{}
	for _, imp := range file.AST.Imports {
		name := r.getImportName(imp)
		if name != "" && name != "." && name != "_" {
			importNames[name] = struct{}{}
		}
	}

	if len(importNames) == 0 {
		return nil
	}

	w := &lintImportShadow{
		importNames: importNames,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoImportShadowRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	importNames := map[string]struct{}{}
	for _, imp := range file.AST.Imports {
		name := r.getImportName(imp)
		if name != "" && name != "." && name != "_" {
			importNames[name] = struct{}{}
		}
	}

	if len(importNames) == 0 {
		return nil
	}

	w := &lintImportShadow{
		importNames: importNames,
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoImportShadowRule) Name() string {
	return "noImportShadow"
}

// Group returns the rule group.
func (*NoImportShadowRule) Group() string {
	return "suspicious"
}

// CacheTier returns the cache tier for this rule.
func (*NoImportShadowRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// getImportName returns the effective name of an import spec.
func (r *NoImportShadowRule) getImportName(imp *ast.ImportSpec) string {
	if imp.Name != nil {
		return imp.Name.Name
	}

	path := strings.Trim(imp.Path.Value, `"`)
	parts := strings.Split(path, "/")

	lastSegment := parts[len(parts)-1]
	if r.isVersion(lastSegment) && len(parts) >= 2 {
		return parts[len(parts)-2]
	}

	return lastSegment
}

// isVersion checks if a path segment looks like a version (v1, v2, etc.).
func (*NoImportShadowRule) isVersion(name string) bool {
	if len(name) < 2 || (name[0] != 'v' && name[0] != 'V') {
		return false
	}
	for i := 1; i < len(name); i++ {
		if name[i] < '0' || name[i] > '9' {
			return false
		}
	}
	return true
}

type lintImportShadow struct {
	importNames map[string]struct{}
	onFailure   func(lint.Failure)
}

func (w *lintImportShadow) Visit(node ast.Node) ast.Visitor {
	assign, ok := node.(*ast.AssignStmt)
	if !ok {
		return w
	}

	// Only check short variable declarations (:=) and plain assignments (=)
	if assign.Tok != token.DEFINE && assign.Tok != token.ASSIGN {
		return w
	}

	for _, lhs := range assign.Lhs {
		ident, ok := lhs.(*ast.Ident)
		if !ok {
			continue
		}

		if _, isImport := w.importNames[ident.Name]; isImport {
			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       ident,
				Category:   lint.FailureCategoryNaming,
				Failure:    fmt.Sprintf("assignment shadows import %q", ident.Name),
			})
		}
	}

	return w
}
