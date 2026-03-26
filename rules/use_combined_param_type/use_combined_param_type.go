package use_combined_param_type

import (
	"fmt"
	"go/ast"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseCombinedParamTypeRule detects if function parameters could be combined by type.
type UseCombinedParamTypeRule struct{}

// Apply applies the rule to given file.
func (r *UseCombinedParamTypeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintCombinedParamType{
		onFailure: onFailure,
		file:      file,
	}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseCombinedParamTypeRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure
	onFailure := func(failure lint.Failure) {
		failures = append(failures, failure)
	}
	w := &lintCombinedParamType{
		onFailure: onFailure,
		file:      file,
	}
	w.Visit(node)
	return failures
}

// Name returns the rule name.
func (*UseCombinedParamTypeRule) Name() string {
	return "useCombinedParamType"
}

// Group returns the rule group.
func (*UseCombinedParamTypeRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseCombinedParamTypeRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintCombinedParamType struct {
	onFailure func(lint.Failure)
	file      *lint.File
}

func (w *lintCombinedParamType) Visit(node ast.Node) ast.Visitor {
	switch n := node.(type) {
	case *ast.FuncDecl:
		w.checkParams(n.Type.Params)
	}
	return w
}

func (w *lintCombinedParamType) checkParams(params *ast.FieldList) {
	if params == nil || len(params.List) < 2 {
		return
	}

	// Skip if any field has no names (unnamed parameter list)
	for _, field := range params.List {
		if len(field.Names) == 0 {
			return
		}
	}

	// Skip multi-line parameter declarations:
	// Check if opening and closing parens are on different lines
	openLine := w.file.ToPosition(params.Opening).Line
	closeLine := w.file.ToPosition(params.Closing).Line
	if openLine != closeLine {
		return
	}

	// Look for consecutive fields with the same type that could be combined
	for i := 0; i < len(params.List)-1; i++ {
		currentType := astutils.GoFmt(params.List[i].Type)
		nextType := astutils.GoFmt(params.List[i+1].Type)
		if currentType == nextType {
			// Collect all consecutive fields with the same type
			var names []string
			j := i
			for j < len(params.List) && astutils.GoFmt(params.List[j].Type) == currentType {
				for _, name := range params.List[j].Names {
					names = append(names, name.Name)
				}
				j++
			}

			w.onFailure(lint.Failure{
				Confidence: 1,
				Node:       params.List[i],
				Category:   lint.FailureCategoryStyle,
				Failure:    fmt.Sprintf("parameters could be combined by type, e.g. %s %s", buildNameList(names), currentType),
			})

			// Skip to the end of this group
			i = j - 1
		}
	}
}

func buildNameList(names []string) string {
	result := ""
	for i, name := range names {
		if i > 0 {
			result += ", "
		}
		result += name
	}
	return result
}
