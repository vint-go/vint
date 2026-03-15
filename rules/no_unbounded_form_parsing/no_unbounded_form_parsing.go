package no_unbounded_form_parsing

import (
	"go/ast"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnboundedFormParsingRule detects unbounded form parsing in HTTP handlers
// that can cause memory exhaustion.
type NoUnboundedFormParsingRule struct{}

// formParsingMethods lists the method names on *http.Request that perform form parsing.
var formParsingMethods = map[string]bool{
	"ParseForm":          true,
	"ParseMultipartForm": true,
	"FormValue":          true,
}

// Apply applies the rule to given file.
func (r *NoUnboundedFormParsingRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}
		checkFunctionForUnboundedFormParsing(funcDecl, onFailure)
	}

	return failures
}

// Name returns the rule name.
func (*NoUnboundedFormParsingRule) Name() string {
	return "noUnboundedFormParsing"
}

// Group returns the rule group.
func (*NoUnboundedFormParsingRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnboundedFormParsingRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// checkFunctionForUnboundedFormParsing checks a function declaration for form
// parsing calls on *http.Request parameters without a preceding http.MaxBytesReader
// call to limit the request body size.
func checkFunctionForUnboundedFormParsing(funcDecl *ast.FuncDecl, onFailure func(lint.Failure)) {
	// Phase 1: Find parameter names that look like *http.Request.
	requestParams := collectHTTPRequestParams(funcDecl)
	if len(requestParams) == 0 {
		return
	}

	// Phase 2: Check if MaxBytesReader is called on the request body anywhere
	// in the function. If so, the request is bounded.
	boundedParams := collectBoundedParams(funcDecl.Body, requestParams)

	// Phase 3: Find form parsing method calls on unbounded request params.
	ast.Inspect(funcDecl.Body, func(n ast.Node) bool {
		ce, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := ce.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}

		if !formParsingMethods[sel.Sel.Name] {
			return true
		}

		ident, ok := sel.X.(*ast.Ident)
		if !ok {
			return true
		}

		if !requestParams[ident.Name] {
			return true
		}

		if boundedParams[ident.Name] {
			return true
		}

		onFailure(lint.Failure{
			Confidence: 1,
			Node:       ce,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "form parsing without http.MaxBytesReader allows unbounded memory consumption",
		})

		return true
	})
}

// collectHTTPRequestParams extracts parameter names whose type is *http.Request.
// It recognizes the pattern: func handler(w http.ResponseWriter, r *http.Request).
func collectHTTPRequestParams(funcDecl *ast.FuncDecl) map[string]bool {
	params := map[string]bool{}
	if funcDecl.Type == nil || funcDecl.Type.Params == nil {
		return params
	}

	for _, field := range funcDecl.Type.Params.List {
		if !isHTTPRequestPointerType(field.Type) {
			continue
		}
		for _, name := range field.Names {
			params[name.Name] = true
		}
	}

	return params
}

// isHTTPRequestPointerType checks if a type expression represents *http.Request.
func isHTTPRequestPointerType(expr ast.Expr) bool {
	star, ok := expr.(*ast.StarExpr)
	if !ok {
		return false
	}
	sel, ok := star.X.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "http" && sel.Sel.Name == "Request"
}

// collectBoundedParams finds request parameters that are bounded by
// http.MaxBytesReader. It looks for patterns like:
//
//	r.Body = http.MaxBytesReader(w, r.Body, limit)
func collectBoundedParams(body *ast.BlockStmt, requestParams map[string]bool) map[string]bool {
	bounded := map[string]bool{}

	ast.Inspect(body, func(n ast.Node) bool {
		assign, ok := n.(*ast.AssignStmt)
		if !ok {
			return true
		}

		// Look for r.Body = http.MaxBytesReader(...)
		for _, rhs := range assign.Rhs {
			if isMaxBytesReaderCall(rhs) {
				// Check if any LHS is r.Body where r is a request param
				for _, lhs := range assign.Lhs {
					sel, ok := lhs.(*ast.SelectorExpr)
					if !ok || sel.Sel.Name != "Body" {
						continue
					}
					ident, ok := sel.X.(*ast.Ident)
					if !ok {
						continue
					}
					if requestParams[ident.Name] {
						bounded[ident.Name] = true
					}
				}

				// Also check arguments of MaxBytesReader for the request param.
				// http.MaxBytesReader(w, r.Body, maxSize) - the second argument
				// tells us which request is being bounded.
				ce := rhs.(*ast.CallExpr)
				if len(ce.Args) >= 2 {
					if bodySel, ok := ce.Args[1].(*ast.SelectorExpr); ok {
						if bodyIdent, ok := bodySel.X.(*ast.Ident); ok {
							if requestParams[bodyIdent.Name] {
								bounded[bodyIdent.Name] = true
							}
						}
					}
				}
			}
		}

		return true
	})

	return bounded
}

// isMaxBytesReaderCall checks if an expression is a call to http.MaxBytesReader.
func isMaxBytesReaderCall(expr ast.Expr) bool {
	ce, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}
	sel, ok := ce.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	ident, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return ident.Name == "http" && sel.Sel.Name == "MaxBytesReader"
}
