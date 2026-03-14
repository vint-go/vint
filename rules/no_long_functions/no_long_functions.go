package no_long_functions

import (
	"fmt"
	"go/ast"
	"strings"

	"github.com/strowk/vint/lint"
)

// NoLongFunctionsRule checks that functions do not exceed a maximum number of lines.
type NoLongFunctionsRule struct {
	maxLines       int
	ignoreComments bool
}

const defaultMaxFunctionLines = 60

// Configure validates the rule configuration, and configures the rule accordingly.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *NoLongFunctionsRule) Configure(arguments lint.Arguments) error {
	if len(arguments) < 1 {
		r.maxLines = defaultMaxFunctionLines
		r.ignoreComments = false
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		// try direct int64 argument for simple threshold configuration
		lines, ok := arguments[0].(int64)
		if !ok {
			return fmt.Errorf(`invalid argument to the "noLongFunctions" rule, expecting a k,v map or int64, got %T`, arguments[0])
		}
		r.maxLines = int(lines)
		return nil
	}

	r.maxLines = defaultMaxFunctionLines
	r.ignoreComments = false
	for k, v := range argKV {
		switch {
		case isRuleOption(k, "lines"):
			lines, ok := v.(int64)
			if !ok {
				return fmt.Errorf(`invalid configuration value for lines in "noLongFunctions" rule; need int64 but got %T`, v)
			}
			r.maxLines = int(lines)
		case isRuleOption(k, "ignoreComments"):
			ignore, ok := v.(bool)
			if !ok {
				return fmt.Errorf(`invalid configuration value for ignoreComments in "noLongFunctions" rule; need bool but got %T`, v)
			}
			r.ignoreComments = ignore
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoLongFunctionsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	if r.maxLines < 0 {
		return nil // disabled when set to -1
	}

	var failures []lint.Failure

	for _, decl := range file.AST.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		body := funcDecl.Body
		if body == nil || len(body.List) == 0 {
			continue
		}

		bodyStartLine := file.ToPosition(body.Lbrace).Line
		bodyEndLine := file.ToPosition(body.Rbrace).Line

		// Lines counted from opening brace to closing brace, excluding the braces themselves
		lineCount := bodyEndLine - bodyStartLine - 1
		if lineCount <= 0 {
			continue
		}

		if r.ignoreComments {
			commentLines := countCommentLinesInBody(file, file.AST.Comments, bodyStartLine, bodyEndLine)
			lineCount -= commentLines
		}

		if lineCount > r.maxLines {
			failures = append(failures, lint.Failure{
				Confidence: 1,
				Category:   lint.FailureCategoryComplexity,
				Failure:    fmt.Sprintf("function %s has too many lines (%d > %d)", funcName(funcDecl), lineCount, r.maxLines),
				Node:       funcDecl,
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoLongFunctionsRule) Name() string {
	return "noLongFunctions"
}

// Group returns the rule group.
func (*NoLongFunctionsRule) Group() string {
	return "complexity"
}

// countCommentLinesInBody counts the number of lines occupied by comments
// within the function body (between startLine and endLine, exclusive).
func countCommentLinesInBody(file *lint.File, comments []*ast.CommentGroup, startLine, endLine int) int {
	count := 0
	for _, cg := range comments {
		for _, comment := range cg.List {
			commentLine := file.ToPosition(comment.Pos()).Line
			if commentLine <= startLine || commentLine >= endLine {
				continue
			}

			if len(comment.Text) < 2 {
				continue
			}
			switch comment.Text[1] {
			case '/': // single-line comment (// ...)
				count++
			case '*': // multi-line comment (/* ... */)
				commentEndLine := file.ToPosition(comment.End()).Line
				// Count only lines within the body range
				firstLine := commentLine
				lastLine := commentEndLine
				if firstLine < startLine+1 {
					firstLine = startLine + 1
				}
				if lastLine > endLine-1 {
					lastLine = endLine - 1
				}
				if lastLine >= firstLine {
					count += lastLine - firstLine + 1
				}
			}
		}
	}
	return count
}

// funcName returns the name representation of a function or method:
// "(Type).Name" for methods or simply "Name" for functions.
func funcName(fn *ast.FuncDecl) string {
	declarationHasReceiver := fn.Recv != nil && fn.Recv.NumFields() > 0
	if declarationHasReceiver {
		typ := fn.Recv.List[0].Type
		return fmt.Sprintf("(%s).%s", recvString(typ), fn.Name)
	}

	return fn.Name.Name
}

// recvString returns a string representation of recv of the
// form "T", "*T", or "BADRECV" (if not a proper receiver type).
func recvString(recv ast.Expr) string {
	switch t := recv.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + recvString(t.X)
	}
	return "BADRECV"
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
