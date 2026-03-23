package no_unescaped_html_template

import (
	"go/ast"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoUnescapedHtmlTemplateRule detects the use of unescaped data in HTML templates,
// which can lead to cross-site scripting (XSS) vulnerabilities.
type NoUnescapedHtmlTemplateRule struct{}

// unsafeHTMLTemplateTypes lists the html/template types that bypass contextual escaping.
var unsafeHTMLTemplateTypes = map[string]bool{
	"HTML": true,
	"JS":   true,
	"URL":  true,
	"CSS":  true,
}

// Apply applies the rule to given file.
func (r *NoUnescapedHtmlTemplateRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if !importsHTMLTemplate(file.AST) {
		return nil
	}

	var failures []lint.Failure

	// Resolve the local name for the html/template package (handles aliases).
	templateAlias := resolveHTMLTemplateAlias(file.AST)

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintUnescapedHTMLTemplate{
		onFailure:     onFailure,
		templateAlias: templateAlias,
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoUnescapedHtmlTemplateRule) Name() string {
	return "noUnescapedHtmlTemplate"
}

// Group returns the rule group.
func (*NoUnescapedHtmlTemplateRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnescapedHtmlTemplateRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintUnescapedHTMLTemplate struct {
	onFailure     func(lint.Failure)
	templateAlias string
}

func (w *lintUnescapedHTMLTemplate) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	// Check if this is a type conversion like template.HTML(expr).
	// In Go AST, type conversions are represented as CallExpr with a
	// SelectorExpr as Fun (e.g., template.HTML is a selector).
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	// Check that the package identifier matches the html/template import.
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return w
	}

	if pkg.Name != w.templateAlias {
		return w
	}

	// Check if the selector is one of the unsafe types.
	if !unsafeHTMLTemplateTypes[sel.Sel.Name] {
		return w
	}

	// Must have exactly one argument (type conversion).
	if len(call.Args) != 1 {
		return w
	}

	// Check if the argument is a variable (not a string literal).
	// String literals are safe because they are developer-controlled constants.
	if isConstantExpr(call.Args[0]) {
		return w
	}

	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       call,
		Category:   lint.FailureCategoryBadPractice,
		Failure:    "use of unescaped data in html/template type conversion template." + sel.Sel.Name,
	})

	return w
}

// isConstantExpr checks if an expression is a compile-time constant
// (string literal or basic literal). Variable references and function
// calls are not considered constant.
func isConstantExpr(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.BasicLit:
		return true
	case *ast.BinaryExpr:
		// Constant concatenation like "a" + "b"
		return isConstantExpr(e.X) && isConstantExpr(e.Y)
	case *ast.ParenExpr:
		return isConstantExpr(e.X)
	}
	return false
}

// importsHTMLTemplate checks if the file imports "html/template".
func importsHTMLTemplate(file *ast.File) bool {
	for _, imp := range file.Imports {
		if imp.Path.Value == `"html/template"` {
			return true
		}
	}
	return false
}

// resolveHTMLTemplateAlias returns the local name used for the html/template package.
// If the import has an alias, that alias is returned; otherwise "template" is returned.
func resolveHTMLTemplateAlias(file *ast.File) string {
	for _, imp := range file.Imports {
		if imp.Path.Value == `"html/template"` {
			if imp.Name != nil {
				return imp.Name.Name
			}
			return "template"
		}
	}
	return "template"
}
