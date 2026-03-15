package no_unsafe_deserialization

import (
	"go/ast"
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoUnsafeDeserializationRule detects unsafe deserialization of untrusted data
// via taint analysis. It traces the flow of untrusted data from sources to
// deserialization sinks such as encoding/gob and encoding/xml.
type NoUnsafeDeserializationRule struct{}

// Apply applies the rule to the given file.
func (r *NoUnsafeDeserializationRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Check if any unsafe deserialization packages are imported.
	hasGob := importsPackage(file.AST, "encoding/gob")
	hasXML := importsPackage(file.AST, "encoding/xml")
	if !hasGob && !hasXML {
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
		if !isHTTPHandler(fn) && fn.Type.Params != nil {
			for _, field := range fn.Type.Params.List {
				if isUntrustedParamType(field.Type) {
					for _, name := range field.Names {
						tainted[name.Name] = true
					}
				}
			}
			findTaintedVars(fn.Body, tainted, funcDecls)
		}

		// Phase 4: Check for deserialization sinks within this function.
		checkDeserializationSinks(fn.Body, tainted, funcDecls, onFailure, hasGob, hasXML)

		// Phase 5: Interprocedural analysis.
		checkInterproceduralSinks(fn.Body, tainted, funcDecls, onFailure, hasGob, hasXML)
	}

	return failures
}

// Name returns the rule name.
func (*NoUnsafeDeserializationRule) Name() string {
	return "noUnsafeDeserialization"
}

// Group returns the rule group.
func (*NoUnsafeDeserializationRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoUnsafeDeserializationRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// importsPackage checks if the file imports a given package path.
func importsPackage(file *ast.File, pkgPath string) bool {
	for _, imp := range file.Imports {
		if imp.Path.Value == `"`+pkgPath+`"` {
			return true
		}
	}
	return false
}

// isUntrustedParamType checks if a type expression represents a type that
// could carry untrusted data (io.Reader, io.ReadCloser, []byte, string).
func isUntrustedParamType(expr ast.Expr) bool {
	switch e := expr.(type) {
	case *ast.Ident:
		return e.Name == "string"
	case *ast.ArrayType:
		if ident, ok := e.Elt.(*ast.Ident); ok {
			return ident.Name == "byte"
		}
	case *ast.SelectorExpr:
		if pkg, ok := e.X.(*ast.Ident); ok {
			if pkg.Name == "io" {
				return e.Sel.Name == "Reader" || e.Sel.Name == "ReadCloser"
			}
		}
	case *ast.StarExpr:
		return isUntrustedParamType(e.X)
	}
	return false
}

// checkDeserializationSinks walks a function body looking for tainted data
// flowing into deserialization calls.
func checkDeserializationSinks(body *ast.BlockStmt, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl, onFailure func(ast.Node, string), hasGob, hasXML bool) {
	if body == nil {
		return
	}

	ast.Inspect(body, func(node ast.Node) bool {
		call, ok := node.(*ast.CallExpr)
		if !ok {
			return true
		}

		if msg := isDeserializationSink(call, tainted, funcDecls, hasGob, hasXML); msg != "" {
			onFailure(call, msg)
			return false
		}

		return true
	})
}

// isDeserializationSink checks if a call expression is a deserialization sink
// with tainted input. Returns the failure message or empty string.
func isDeserializationSink(call *ast.CallExpr, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl, hasGob, hasXML bool) string {
	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return ""
	}

	// Check for gob.NewDecoder(taintedReader) pattern
	if hasGob {
		if pkg, ok := sel.X.(*ast.Ident); ok {
			if pkg.Name == "gob" && sel.Sel.Name == "NewDecoder" && len(call.Args) == 1 {
				if isTaintedExpr(call.Args[0], tainted, funcDecls) {
					return "unsafe deserialization: untrusted data passed to gob.NewDecoder"
				}
			}
		}
	}

	// Check for decoder.Decode() where the decoder was created from tainted source
	if sel.Sel.Name == "Decode" && len(call.Args) == 1 {
		if ident, ok := sel.X.(*ast.Ident); ok {
			if tainted[ident.Name] {
				return "unsafe deserialization: calling Decode on a decoder created from untrusted source"
			}
		}
	}

	// Check for xml.Unmarshal(taintedData, ...) pattern
	if hasXML {
		if pkg, ok := sel.X.(*ast.Ident); ok {
			if pkg.Name == "xml" && sel.Sel.Name == "Unmarshal" && len(call.Args) == 2 {
				if isTaintedExpr(call.Args[0], tainted, funcDecls) {
					return "unsafe deserialization: untrusted data passed to xml.Unmarshal"
				}
			}
		}

		// Check for xml.NewDecoder(taintedReader) pattern
		if pkg, ok := sel.X.(*ast.Ident); ok {
			if pkg.Name == "xml" && sel.Sel.Name == "NewDecoder" && len(call.Args) == 1 {
				if isTaintedExpr(call.Args[0], tainted, funcDecls) {
					return "unsafe deserialization: untrusted data passed to xml.NewDecoder"
				}
			}
		}
	}

	return ""
}

// checkInterproceduralSinks checks calls to local functions that receive
// tainted args and contain deserialization sinks internally.
func checkInterproceduralSinks(body *ast.BlockStmt, tainted map[string]bool, funcDecls map[string]*ast.FuncDecl, onFailure func(ast.Node, string), hasGob, hasXML bool) {
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

		// Check if callee contains deserialization sinks with tainted data.
		checkDeserializationSinks(fn.Body, calleeTainted, funcDecls, func(_ ast.Node, msg string) {
			onFailure(call, msg)
		}, hasGob, hasXML)

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

	// io.ReadAll(r.Body) or ioutil.ReadAll(r.Body) pattern
	if pkg, ok := sel.X.(*ast.Ident); ok {
		if (pkg.Name == "io" || pkg.Name == "ioutil") && sel.Sel.Name == "ReadAll" {
			if len(call.Args) == 1 {
				if isTaintedReader(call.Args[0]) {
					return true
				}
			}
		}
	}

	return false
}

// isTaintedReader checks if an expression is a reader from an HTTP request body.
func isTaintedReader(expr ast.Expr) bool {
	sel, ok := expr.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	// r.Body
	return sel.Sel.Name == "Body"
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
		"io":       {"ReadAll": true},
		"ioutil":   {"ReadAll": true},
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
