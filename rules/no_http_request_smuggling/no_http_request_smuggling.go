package no_http_request_smuggling

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoHttpRequestSmugglingRule detects potential HTTP request smuggling vulnerabilities
// caused by conflicting headers or bare line feed (LF) characters in HTTP requests.
type NoHttpRequestSmugglingRule struct{}

// Apply applies the rule to given file.
func (r *NoHttpRequestSmugglingRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Check for bare LF in HTTP request strings
	bareLFWalker := &lintBareLFInHTTP{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(bareLFWalker, file.AST)

	// Check for conflicting headers within functions
	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok || funcDecl.Body == nil {
			continue
		}
		checkConflictingHeaders(funcDecl.Body, func(f lint.Failure) {
			failures = append(failures, f)
		})
	}

	return failures
}

// Name returns the rule name.
func (*NoHttpRequestSmugglingRule) Name() string {
	return "noHttpRequestSmuggling"
}

// Group returns the rule group.
func (*NoHttpRequestSmugglingRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoHttpRequestSmugglingRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// lintBareLFInHTTP walks AST looking for string literals that contain
// HTTP protocol patterns with bare LF characters (without CR).
type lintBareLFInHTTP struct {
	onFailure func(lint.Failure)
}

// httpIndicators are substrings that indicate an HTTP-like request string.
var httpIndicators = []string{
	"HTTP/1.0",
	"HTTP/1.1",
	"HTTP/2",
	"Content-Length:",
	"Transfer-Encoding:",
	"Host:",
}

func (w *lintBareLFInHTTP) Visit(node ast.Node) ast.Visitor {
	lit, ok := node.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return w
	}

	val, err := strconv.Unquote(lit.Value)
	if err != nil {
		return w
	}

	if containsBareLFInHTTPContext(val) {
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       lit,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    "potential HTTP request smuggling: use CRLF (\\r\\n) instead of bare LF (\\n) in HTTP requests",
		})
	}

	return w
}

// containsBareLFInHTTPContext checks if a string contains bare LF characters
// in what appears to be an HTTP request context.
func containsBareLFInHTTPContext(s string) bool {
	// Must contain a bare LF (i.e., \n not preceded by \r)
	if !containsBareLF(s) {
		return false
	}

	// Must look like an HTTP request or header
	for _, indicator := range httpIndicators {
		if strings.Contains(s, indicator) {
			return true
		}
	}

	return false
}

// containsBareLF returns true if the string contains at least one \n
// that is not preceded by \r.
func containsBareLF(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' && (i == 0 || s[i-1] != '\r') {
			return true
		}
	}
	return false
}

// headerSetInfo records information about a Header.Set call.
type headerSetInfo struct {
	headerName string
	node       ast.Node
	varName    string
}

// checkConflictingHeaders inspects a function body for cases where both
// Content-Length and Transfer-Encoding headers are set on the same request variable.
func checkConflictingHeaders(body *ast.BlockStmt, onFailure func(lint.Failure)) {
	var headerSets []headerSetInfo

	ast.Inspect(body, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok || sel.Sel.Name != "Set" {
			return true
		}

		// Check for pattern: <expr>.Header.Set(<headerName>, <value>)
		innerSel, ok := sel.X.(*ast.SelectorExpr)
		if !ok || innerSel.Sel.Name != "Header" {
			return true
		}

		// Get the variable name the Header belongs to
		varIdent, ok := innerSel.X.(*ast.Ident)
		if !ok {
			return true
		}

		// Get the header name from the first argument
		if len(call.Args) < 1 {
			return true
		}
		headerLit, ok := call.Args[0].(*ast.BasicLit)
		if !ok || headerLit.Kind != token.STRING {
			return true
		}
		headerName, err := strconv.Unquote(headerLit.Value)
		if err != nil {
			return true
		}

		headerSets = append(headerSets, headerSetInfo{
			headerName: strings.ToLower(headerName),
			node:       call,
			varName:    varIdent.Name,
		})

		return true
	})

	// Check for conflicting headers on the same variable
	for i, h1 := range headerSets {
		for _, h2 := range headerSets[i+1:] {
			if h1.varName != h2.varName {
				continue
			}
			if (h1.headerName == "content-length" && h2.headerName == "transfer-encoding") ||
				(h1.headerName == "transfer-encoding" && h2.headerName == "content-length") {
				// Report on the second header set (the one that creates the conflict)
				onFailure(lint.Failure{
					Confidence: 1,
					Node:       h2.node,
					Category:   lint.FailureCategoryBadPractice,
					Failure:    "potential HTTP request smuggling: conflicting Content-Length and Transfer-Encoding headers",
				})
			}
		}
	}
}
