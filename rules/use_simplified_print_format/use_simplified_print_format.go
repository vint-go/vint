package use_simplified_print_format

import (
	"go/ast"

	"github.com/strowk/vint/internal/astutils"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// UseSimplifiedPrintFormatRule detects unnecessarily complex ways of printing
// formatted strings, such as fmt.Print(fmt.Sprintf(...)) which can be simplified
// to fmt.Printf(...).
type UseSimplifiedPrintFormatRule struct{}

// Apply applies the rule to the given file.
func (r *UseSimplifiedPrintFormatRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSimplifiedPrint{onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseSimplifiedPrintFormatRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}

	w := &lintSimplifiedPrint{onFailure: onFailure}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseSimplifiedPrintFormatRule) Name() string {
	return "useSimplifiedPrintFormat"
}

// Group returns the rule group.
func (*UseSimplifiedPrintFormatRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseSimplifiedPrintFormatRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintSimplifiedPrint struct {
	onFailure func(lint.Failure)
}

// simplification maps outer print functions to their simplified form.
// The key is the outer function name, and the value contains the replacement
// and the expected argument count (how many args the outer call should have).
type printMapping struct {
	replacement string
	// argCount is the total number of arguments expected in the outer call.
	// For Print/Println it's 1 (just the Sprintf result).
	// For Fprint/Fprintln it's 2 (the writer + the Sprintf result).
	argCount int
	// sprintfArgIndex is the index of the argument that should be fmt.Sprintf.
	sprintfArgIndex int
}

var printMappings = map[string]printMapping{
	"Print":   {replacement: "fmt.Printf", argCount: 1, sprintfArgIndex: 0},
	"Println": {replacement: "fmt.Printf", argCount: 1, sprintfArgIndex: 0},
	"Fprint":  {replacement: "fmt.Fprintf", argCount: 2, sprintfArgIndex: 1},
	"Fprintln": {replacement: "fmt.Fprintf", argCount: 2, sprintfArgIndex: 1},
}

func (w *lintSimplifiedPrint) Visit(node ast.Node) ast.Visitor {
	call, ok := node.(*ast.CallExpr)
	if !ok {
		return w
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return w
	}

	if !astutils.IsIdent(sel.X, "fmt") {
		return w
	}

	mapping, ok := printMappings[sel.Sel.Name]
	if !ok {
		return w
	}

	if len(call.Args) != mapping.argCount {
		return w
	}

	// Check if the relevant argument is a fmt.Sprintf call
	innerCall, ok := call.Args[mapping.sprintfArgIndex].(*ast.CallExpr)
	if !ok {
		return w
	}

	if !astutils.IsPkgDotName(innerCall.Fun, "fmt", "Sprintf") {
		return w
	}

	outerFunc := "fmt." + sel.Sel.Name
	w.onFailure(lint.Failure{
		Confidence: 1,
		Node:       call,
		Category:   lint.FailureCategoryStyle,
		Failure:    outerFunc + "(fmt.Sprintf(...)) can be simplified to " + mapping.replacement + "(...)",
	})

	return w
}
