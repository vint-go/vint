package no_ssrf_taint

import (
	"go/ast"
	"go/token"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoSsrfTaintRule detects Server-Side Request Forgery (SSRF) vulnerabilities
// via taint analysis. It traces the flow of untrusted data from sources (such
// as HTTP request parameters) to HTTP request sinks (such as http.Get,
// http.Post, or http.NewRequest). It performs interprocedural analysis to
// detect cases where user-controlled URLs or URL components flow through
// multiple functions before being used in outbound HTTP requests.
type NoSsrfTaintRule struct{}

// Apply applies the rule to the given file.
func (r *NoSsrfTaintRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Collect all function declarations for interprocedural analysis.
	funcDecls := map[string]*ast.FuncDecl{}
	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}
		funcDecls[fn.Name.Name] = fn
	}

	onFailure := func(node ast.Node, msg string) {
		failures = append(failures, lint.Failure{
			Confidence: 1,
			Node:       node,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    msg,
		})
	}

	for _, decl := range file.AST.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}

		tainted := map[string]bool{}

		// Phase 1: Seed taint from direct taint sources within this function.
		findTaintedVars(fn.Body, tainted, funcDecls)

		// Phase 2: If this is an HTTP handler, the *http.Request param is tainted.
		if isHTTPHandler(fn) {
			if reqParam := findHTTPRequestParam(fn); reqParam != "" {
				tainted[reqParam] = true
				// Re-propagate taint with the seeded request param.
				findTaintedVars(fn.Body, tainted, funcDecls)
			}
		}

		// Phase 3: Check for direct sinks within this function.
		checkSinks(fn.Body, tainted, funcDecls, onFailure)

		// Phase 4: Interprocedural - check calls to local functions with tainted args.
		checkInterproceduralSinks(fn.Body, tainted, funcDecls, onFailure)
	}

	return failures
}

// Name returns the rule name.
func (*NoSsrfTaintRule) Name() string {
	return "noSsrfTaint"
}

// Group returns the rule group.
func (*NoSsrfTaintRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoSsrfTaintRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// httpSinkFunctions maps net/http package-level functions to the
// zero-based index of their URL argument.
var httpSinkFunctions = map[string]int{
	"Get":      0, // http.Get(url)
	"Head":     0, // http.Head(url)
	"Post":     0, // http.Post(url, contentType, body)
	"PostForm": 0, // http.PostForm(url, data)
}

// isTaintSourceCall checks if a call expression is a taint source.
// Returns true if the expression produces user-controlled data.
func isTaintSourceCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// r.URL.Query().Get("...") pattern: the outer call is to .Get()
	if sel.Sel.Name == "Get" {
		if innerCall, ok := sel.X.(*ast.CallExpr); ok {
			if isURLQueryCall(innerCall) {
				return true
			}
		}
	}

	// r.FormValue("..."), r.PostFormValue("...")
	switch sel.Sel.Name {
	case "FormValue", "PostFormValue":
		return true
	}

	// r.URL.Query() itself returns tainted values
	if isURLQueryCall(call) {
		return true
	}

	// bufio.NewScanner(os.Stdin), bufio.NewReader(os.Stdin)
	if pkg, ok := sel.X.(*ast.Ident); ok {
		if pkg.Name == "bufio" && (sel.Sel.Name == "NewScanner" || sel.Sel.Name == "NewReader") {
			for _, arg := range call.Args {
				if isOsStdin(arg) {
					return true
				}
			}
		}
	}

	return false
}

// isURLQueryCall checks if a call expression is r.URL.Query().
func isURLQueryCall(call *ast.CallExpr) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Query" {
		return false
	}
	if innerSel, ok := sel.X.(*ast.SelectorExpr); ok {
		return innerSel.Sel.Name == "URL"
	}
	return false
}

// isOsStdin checks if an expression is os.Stdin.
func isOsStdin(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}
	return pkg.Name == "os" && sel.Sel.Name == "Stdin"
}

// findTaintedVars walks a function body to find variables that hold tainted data.
func findTaintedVars(body *ast.BlockStmt, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl) {
	if body == nil {
		return
	}
	for _, stmt := range body.List {
		findTaintedInStmt(stmt, tainted, funcDecls)
	}
}

// findTaintedInStmt processes a single statement for taint tracking.
func findTaintedInStmt(stmt ast.Stmt, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		// For multi-return with single RHS (e.g., out, err := fn(x))
		if len(s.Rhs) == 1 && len(s.Lhs) >= 1 {
			if isTaintedExpr(s.Rhs[0], tainted, funcDecls) {
				for _, lhs := range s.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok && ident.Name != "_" {
						tainted[ident.Name] = true
					}
				}
			}
		} else {
			// Positional matching
			for i, rhs := range s.Rhs {
				if isTaintedExpr(rhs, tainted, funcDecls) && i < len(s.Lhs) {
					if ident, ok := s.Lhs[i].(*ast.Ident); ok {
						tainted[ident.Name] = true
					}
				}
			}
		}
	case *ast.DeclStmt:
		if genDecl, ok := s.Decl.(*ast.GenDecl); ok && genDecl.Tok == token.VAR {
			for _, spec := range genDecl.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					for j, val := range vs.Values {
						if isTaintedExpr(val, tainted, funcDecls) && j < len(vs.Names) {
							tainted[vs.Names[j].Name] = true
						}
					}
				}
			}
		}
	case *ast.IfStmt:
		if s.Init != nil {
			findTaintedInStmt(s.Init, tainted, funcDecls)
		}
		findTaintedVars(s.Body, tainted, funcDecls)
		if s.Else != nil {
			if elseBlock, ok := s.Else.(*ast.BlockStmt); ok {
				findTaintedVars(elseBlock, tainted, funcDecls)
			} else if elseIf, ok := s.Else.(*ast.IfStmt); ok {
				findTaintedInStmt(elseIf, tainted, funcDecls)
			}
		}
	case *ast.ForStmt:
		if s.Init != nil {
			findTaintedInStmt(s.Init, tainted, funcDecls)
		}
		findTaintedVars(s.Body, tainted, funcDecls)
	case *ast.RangeStmt:
		if s.Value != nil {
			if ident, ok := s.Value.(*ast.Ident); ok {
				if rangeExpr, ok := s.X.(*ast.Ident); ok && tainted[rangeExpr.Name] {
					tainted[ident.Name] = true
				}
			}
		}
		findTaintedVars(s.Body, tainted, funcDecls)
	case *ast.BlockStmt:
		findTaintedVars(s, tainted, funcDecls)
	}
}

// isTaintedExpr checks if an expression is tainted.
func isTaintedExpr(expr ast.Expr, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		return tainted[e.Name]
	case *ast.CallExpr:
		if isTaintSourceCall(expr) {
			return true
		}
		// Check if calling a local function that propagates taint through return values.
		if ident, ok := e.Fun.(*ast.Ident); ok {
			if fn, ok := funcDecls[ident.Name]; ok {
				return callPropagatesTaint(e, fn, tainted, funcDecls)
			}
		}
		// Check taint-propagating stdlib calls (fmt.Sprintf, strings.Join, etc.)
		if isTaintPropagatingCall(e, tainted, funcDecls) {
			return true
		}
	case *ast.BinaryExpr:
		return isTaintedExpr(e.X, tainted, funcDecls) || isTaintedExpr(e.Y, tainted, funcDecls)
	case *ast.IndexExpr:
		return isTaintedExpr(e.X, tainted, funcDecls)
	case *ast.SliceExpr:
		return isTaintedExpr(e.X, tainted, funcDecls)
	case *ast.SelectorExpr:
		return isTaintedExpr(e.X, tainted, funcDecls)
	case *ast.ParenExpr:
		return isTaintedExpr(e.X, tainted, funcDecls)
	case *ast.TypeAssertExpr:
		return isTaintedExpr(e.X, tainted, funcDecls)
	case *ast.UnaryExpr:
		return isTaintedExpr(e.X, tainted, funcDecls)
	}
	return false
}

// isTaintPropagatingCall checks if a function call propagates taint from its arguments.
func isTaintPropagatingCall(call *ast.CallExpr, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false
	}

	propagators := map[string]map[string]bool{
		"fmt":      {"Sprintf": true, "Sprint": true, "Sprintln": true},
		"strings":  {"Join": true, "Replace": true, "ToLower": true, "ToUpper": true, "TrimSpace": true, "Trim": true},
		"filepath": {"Join": true},
		"path":     {"Join": true},
		"url":      {"JoinPath": true},
	}

	if methods, ok := propagators[pkg.Name]; ok && methods[sel.Sel.Name] {
		for _, arg := range call.Args {
			if isTaintedExpr(arg, tainted, funcDecls) {
				return true
			}
		}
	}

	return false
}

// callPropagatesTaint checks if calling a local function with tainted arguments
// will produce tainted return values (interprocedural analysis).
func callPropagatesTaint(call *ast.CallExpr, fn *ast.FuncDecl, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl) bool {
	if fn.Body == nil || fn.Type.Params == nil {
		return false
	}

	// Check if any call argument is tainted.
	hasTaintedArg := false
	for _, arg := range call.Args {
		if isTaintedExpr(arg, tainted, funcDecls) {
			hasTaintedArg = true
			break
		}
	}
	if !hasTaintedArg {
		return false
	}

	// Map tainted call arguments to function parameter names.
	paramTainted := buildParamTaintMap(call, fn, tainted, funcDecls)

	// Check if the function returns tainted data.
	return functionReturnsTaint(fn.Body, paramTainted, funcDecls)
}

// buildParamTaintMap maps tainted call arguments to function parameter names.
func buildParamTaintMap(call *ast.CallExpr, fn *ast.FuncDecl, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl) map[string]bool {
	paramTainted := map[string]bool{}
	paramIdx := 0
	for _, field := range fn.Type.Params.List {
		for _, name := range field.Names {
			if paramIdx < len(call.Args) && isTaintedExpr(call.Args[paramIdx], tainted, funcDecls) {
				paramTainted[name.Name] = true
			}
			paramIdx++
		}
	}
	return paramTainted
}

// functionReturnsTaint checks if a function body returns tainted data
// given the set of tainted parameter names.
func functionReturnsTaint(body *ast.BlockStmt, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl) bool {
	localTainted := map[string]bool{}
	for k, v := range tainted {
		localTainted[k] = v
	}
	findTaintedVars(body, localTainted, funcDecls)

	found := false
	ast.Inspect(body, func(node ast.Node) bool {
		if found {
			return false
		}
		ret, ok := node.(*ast.ReturnStmt)
		if !ok {
			return true
		}
		for _, result := range ret.Results {
			if isTaintedExpr(result, localTainted, funcDecls) {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

// isHTTPSink checks if a call expression is an outbound HTTP request sink.
// It returns true for calls like http.Get, http.Post, http.Head, http.PostForm,
// http.NewRequest, and http.NewRequestWithContext.
func isHTTPSink(call *ast.CallExpr) (isSink bool, urlArgIdx int) {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false, -1
	}

	// Check for package-level http.Xxx calls
	pkg, ok := sel.X.(*ast.Ident)
	if !ok {
		return false, -1
	}

	if pkg.Name != "http" {
		return false, -1
	}

	// Check package-level functions: http.Get, http.Head, http.Post, http.PostForm
	if idx, ok := httpSinkFunctions[sel.Sel.Name]; ok {
		return true, idx
	}

	// http.NewRequest(method, url, body) - URL is the second argument (index 1)
	if sel.Sel.Name == "NewRequest" {
		return true, 1
	}

	// http.NewRequestWithContext(ctx, method, url, body) - URL is the third argument (index 2)
	if sel.Sel.Name == "NewRequestWithContext" {
		return true, 2
	}

	return false, -1
}

// checkSinks walks a function body looking for tainted data flowing into HTTP request sinks.
func checkSinks(body *ast.BlockStmt, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl, onFailure func(ast.Node, string)) {
	if body == nil {
		return
	}

	ast.Inspect(body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		isSink, urlArgIdx := isHTTPSink(call)
		if !isSink {
			return true
		}

		if len(call.Args) <= urlArgIdx {
			return true
		}

		urlArg := call.Args[urlArgIdx]

		if isTaintedExpr(urlArg, tainted, funcDecls) {
			onFailure(call, "potential SSRF: tainted data from user input flows to outbound HTTP request")
		}

		return true
	})
}

// checkInterproceduralSinks checks calls to local functions that receive tainted args
// and contain HTTP request sinks internally.
func checkInterproceduralSinks(body *ast.BlockStmt, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl, onFailure func(ast.Node, string)) {
	if body == nil {
		return
	}

	ast.Inspect(body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		ident, ok := call.Fun.(*ast.Ident)
		if !ok {
			return true
		}

		fn, ok := funcDecls[ident.Name]
		if !ok || fn.Body == nil {
			return true
		}

		// Check if any argument is tainted.
		hasTaintedArg := false
		for _, arg := range call.Args {
			if isTaintedExpr(arg, tainted, funcDecls) {
				hasTaintedArg = true
				break
			}
		}
		if !hasTaintedArg {
			return true
		}

		// Map tainted arguments to parameters.
		calleeTainted := buildParamTaintMap(call, fn, tainted, funcDecls)

		// Propagate taint within the callee.
		findTaintedVars(fn.Body, calleeTainted, funcDecls)

		// Check if callee contains HTTP request sinks with tainted data.
		checkSinks(fn.Body, calleeTainted, funcDecls, func(_ ast.Node, msg string) {
			// Report the failure at the call site in the caller.
			onFailure(call, msg)
		})

		return true
	})
}

// isHTTPHandler checks if a function declaration looks like an HTTP handler.
func isHTTPHandler(fn *ast.FuncDecl) bool {
	if fn.Type.Params == nil || len(fn.Type.Params.List) < 2 {
		return false
	}

	hasResponseWriter := false
	hasRequest := false

	for _, field := range fn.Type.Params.List {
		typeStr := exprToString(field.Type)
		if typeStr == "http.ResponseWriter" {
			hasResponseWriter = true
		}
		if typeStr == "*http.Request" {
			hasRequest = true
		}
	}

	return hasResponseWriter && hasRequest
}

// findHTTPRequestParam returns the name of the *http.Request parameter.
func findHTTPRequestParam(fn *ast.FuncDecl) string {
	for _, field := range fn.Type.Params.List {
		typeStr := exprToString(field.Type)
		if typeStr == "*http.Request" && len(field.Names) > 0 {
			return field.Names[0].Name
		}
	}
	return ""
}

// exprToString returns a simple string representation of a type expression.
func exprToString(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name
	case *ast.SelectorExpr:
		return exprToString(e.X) + "." + e.Sel.Name
	case *ast.StarExpr:
		return "*" + exprToString(e.X)
	}
	return ""
}
