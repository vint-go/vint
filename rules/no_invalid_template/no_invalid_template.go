package no_invalid_template

import (
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"text/template/parse"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoInvalidTemplateRule validates that template strings passed to template.Must,
// template.New().Parse(), and related functions are valid Go templates.
// Invalid templates will cause runtime errors.
type NoInvalidTemplateRule struct{}

// Apply applies the rule to given file.
func (r *NoInvalidTemplateRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintInvalidTemplate{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *NoInvalidTemplateRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintInvalidTemplate{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*NoInvalidTemplateRule) Name() string {
	return "noInvalidTemplate"
}

// Group returns the rule group.
func (*NoInvalidTemplateRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoInvalidTemplateRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintInvalidTemplate struct {
	onFailure func(lint.Failure)
}

func (w *lintInvalidTemplate) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Look for calls that end up passing a string literal to Parse.
	// Patterns we want to detect:
	// 1. template.Must(template.New("name").Parse("..."))
	// 2. template.New("name").Parse("...")
	// 3. t.Parse("...") where t is from text/template or html/template
	//
	// We focus on detecting .Parse("literal") call expressions.
	templateStr, found := extractTemplateString(call)
	if !found {
		return w
	}

	// Try to parse the template string
	_, err := parse.Parse("", templateStr, "", "", nil)
	if err != nil {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryLogic,
			Confidence: 1,
			Node:       call,
			Failure:    fmt.Sprintf("invalid template: %s", err),
		})
	}

	return w
}

// extractTemplateString tries to extract a template string from a call expression.
// It looks for patterns like:
//   - template.New("name").Parse("templateStr")
//   - t.Parse("templateStr")
//   - template.Must(template.New("name").Parse("templateStr")) - the inner Parse call
//
// Returns the template string and true if found, empty string and false otherwise.
func extractTemplateString(call *ast.CallExpr) (string, bool) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", false
	}

	// Check for .Parse("...") calls
	if sel.Sel.Name == "Parse" {
		if len(call.Args) < 1 {
			return "", false
		}
		lit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || lit.Kind != token.STRING {
			return "", false
		}
		s, err := strconv.Unquote(lit.Value)
		if err != nil {
			return "", false
		}

		// Verify the receiver looks like a template:
		// either template.New("...") or an identifier (variable)
		if isTemplateNewCall(sel.X) || isIdent(sel.X) {
			return s, true
		}
		return "", false
	}

	return "", false
}

// isTemplateNewCall checks if expr looks like template.New("name") call.
func isTemplateNewCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return (ident.Name == "template") && sel.Sel.Name == "New"
}

// isIdent checks if expr is a simple identifier (variable name).
func isIdent(expr ast.Expr) bool {
	_, ok := expr.(*ast.Ident)
	return ok
}
