package no_template_injection

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoTemplateInjectionRule detects Server-Side Template Injection (SSTI)
// vulnerabilities via text/template. It traces the flow of untrusted data
// from sources to template parsing and execution sinks.
type NoTemplateInjectionRule struct{}

// Apply applies the rule to the given file.
func (r *NoTemplateInjectionRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Check if text/template is imported
	if !importsTextTemplate(file.AST) {
		return nil
	}

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

		// Phase 3: Mark function parameters as tainted for non-handler functions.
		// Parameters of non-handler helper functions that receive string-like
		// values could carry tainted data.
		if !isHTTPHandler(fn) && fn.Type.Params != nil {
			for _, field := range fn.Type.Params.List {
				if isStringType(field.Type) {
					for _, name := range field.Names {
						tainted[name.Name] = true
					}
				}
			}
			// Re-propagate taint after seeding params.
			findTaintedVars(fn.Body, tainted, funcDecls)
		}

		// Phase 4: Check for template injection sinks within this function.
		checkTemplateSinks(fn.Body, tainted, funcDecls, onFailure)

		// Phase 5: Interprocedural - check calls to local functions with tainted args.
		checkInterproceduralTemplateSinks(fn.Body, tainted, funcDecls, onFailure)
	}

	return failures
}

// Name returns the rule name.
func (*NoTemplateInjectionRule) Name() string {
	return "noTemplateInjection"
}

// Group returns the rule group.
func (*NoTemplateInjectionRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoTemplateInjectionRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// importsTextTemplate checks if the file imports "text/template".
func importsTextTemplate(file *ast.File) bool {
	for _, imp := range file.Imports {
		if imp.Path.Value == `"text/template"` {
			return true
		}
	}
	return false
}

// isStringType checks if a type expression is a string or []byte type.
func isStringType(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name == "string"
	case *ast.ArrayType:
		if ident, ok := e.Elt.(*ast.Ident); ok {
			return ident.Name == "byte"
		}
	}
	return false
}

// isTemplateParseSink checks if a call expression is a template.Parse() or
// template.New(...).Parse() sink that takes tainted input.
// Returns the call and whether it is a sink.
func isTemplateParseSink(call *ast.CallExpr, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl) bool {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// Check for .Parse(taintedArg) method call
	if sel.Sel.Name == "Parse" && len(call.Args) == 1 {
		arg := call.Args[0]
		if isTaintedExpr(arg, tainted, funcDecls) {
			// Now verify the receiver is a text/template object.
			// This covers: template.New("x").Parse(tainted)
			//              t.Parse(tainted) where t is a *template.Template
			if isTemplateReceiver(sel.X, tainted) {
				return true
			}
		}
	}

	return false
}

// isTemplateReceiver checks if an expression is likely a text/template.Template object.
func isTemplateReceiver(expr ast.Expr, tainted map[string]bool) bool {
	switch e := expr.(type) {
	case *ast.CallExpr:
		// template.New("name") or template.Must(...)
		if sel, ok := e.Fun.(*ast.SelectorExpr); ok {
			if pkg, ok := sel.X.(*ast.Ident); ok {
				if pkg.Name == "template" && (sel.Sel.Name == "New" || sel.Sel.Name == "Must") {
					return true
				}
			}
			// Also handle chained calls like t.Funcs(...), t.Option(...), etc.
			if isTemplateReceiver(sel.X, tainted) {
				return true
			}
		}
	case *ast.Ident:
		// A variable that is likely a template — we accept any ident that is not
		// known to be something else. This is conservative but appropriate for
		// security rules.
		return true
	}
	return false
}

// checkTemplateSinks walks a function body looking for tainted data
// flowing into template.Parse() calls.
func checkTemplateSinks(body *ast.BlockStmt, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl, onFailure func(ast.Node, string)) {
	if body == nil {
		return
	}

	ast.Inspect(body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		if isTemplateParseSink(call, tainted, funcDecls) {
			onFailure(call, "potential server-side template injection: tainted data used in text/template.Parse")
			return false
		}

		return true
	})
}

// checkInterproceduralTemplateSinks checks calls to local functions that receive
// tainted args and contain template sinks internally.
func checkInterproceduralTemplateSinks(body *ast.BlockStmt, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl, onFailure func(ast.Node, string)) {
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

		// Check if callee contains template sinks with tainted data.
		checkTemplateSinks(fn.Body, calleeTainted, funcDecls, func(_ ast.Node, msg string) {
			// Report the failure at the call site in the caller.
			onFailure(call, msg)
		})

		return true
	})
}

// isTaintSourceCall checks if a call expression is a taint source.
func isTaintSourceCall(expr ast.Expr) bool {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return false
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	// r.URL.Query().Get("...") pattern
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
		if len(s.Rhs) == 1 && len(s.Lhs) >= 1 {
			if isTaintedExpr(s.Rhs[0], tainted, funcDecls) {
				for _, lhs := range s.Lhs {
					if ident, ok := lhs.(*ast.Ident); ok && ident.Name != "_" {
						tainted[ident.Name] = true
					}
				}
			}
		} else {
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
		if ident, ok := e.Fun.(*ast.Ident); ok {
			if fn, ok := funcDecls[ident.Name]; ok {
				return callPropagatesTaint(e, fn, tainted, funcDecls)
			}
		}
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
// will produce tainted return values.
func callPropagatesTaint(call *ast.CallExpr, fn *ast.FuncDecl, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl) bool {
	if fn.Body == nil || fn.Type.Params == nil {
		return false
	}

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

	paramTainted := buildParamTaintMap(call, fn, tainted, funcDecls)
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

// functionReturnsTaint checks if a function body returns tainted data.
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
