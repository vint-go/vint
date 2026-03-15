package no_commented_out_code

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// minCommentLen is the default minimum length of a comment text (after stripping
// the comment markers) for it to be considered a candidate for commented-out code.
const minCommentLen = 15

// NoCommentedOutCodeRule detects commented-out code inside function bodies.
// It uses heuristics to identify comments that appear to contain executable Go
// code rather than actual documentation.
type NoCommentedOutCodeRule struct{}

// Apply applies the rule to given file.
func (r *NoCommentedOutCodeRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	// Collect the line ranges of all function bodies so we can restrict
	// detection to comments that live inside functions.
	type lineRange struct{ start, end int }
	var funcRanges []lineRange

	ast.Inspect(file.AST, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			return true
		}
		start := file.ToPosition(fn.Body.Lbrace).Line
		end := file.ToPosition(fn.Body.Rbrace).Line
		funcRanges = append(funcRanges, lineRange{start, end})
		return true
	})

	insideFunc := func(line int) bool {
		for _, r := range funcRanges {
			if line >= r.start && line <= r.end {
				return true
			}
		}
		return false
	}

	for _, group := range file.AST.Comments {
		// Check if at least the start of the comment group is inside a function body.
		commentLine := file.ToPosition(group.Pos()).Line
		if !insideFunc(commentLine) {
			continue
		}

		// Merge all lines in the comment group into a single text blob
		// and check for commented-out code.
		for _, comment := range group.List {
			text := commentText(comment.Text)

			if len(text) < minCommentLen {
				continue
			}

			if hasExplanatoryMarker(text) {
				continue
			}

			if looksLikeCode(text) {
				failures = append(failures, lint.Failure{
					Category:   lint.FailureCategoryStyle,
					Confidence: 0.8,
					Node:       comment,
					Failure:    fmt.Sprintf("comment appears to contain commented-out code: %s", truncate(text, 40)),
				})
			}
		}
	}

	return failures
}

// commentText strips the comment markers from a single comment line.
func commentText(raw string) string {
	if strings.HasPrefix(raw, "//") {
		return strings.TrimSpace(raw[2:])
	}
	if strings.HasPrefix(raw, "/*") && strings.HasSuffix(raw, "*/") {
		return strings.TrimSpace(raw[2 : len(raw)-2])
	}
	return raw
}

// hasExplanatoryMarker returns true if the comment text contains markers
// that indicate it is explanatory rather than commented-out code.
func hasExplanatoryMarker(text string) bool {
	lower := strings.ToLower(text)
	// Substring markers that are safe to match anywhere in the text.
	substringMarkers := []string{
		"todo", "fixme", "hack", "xxx",
		"http://", "https://",
		"e.g.", "i.e.", "cf.",
		"nolint", "lint:",
		"deprecated", "warning:",
		"copyright", "license",
	}
	for _, m := range substringMarkers {
		if strings.Contains(lower, m) {
			return true
		}
	}
	// Prefix markers that only match at the start of the comment text.
	prefixMarkers := []string{
		"note:", "note ",
		"bug:", "bug(",
		"see ",
	}
	for _, m := range prefixMarkers {
		if strings.HasPrefix(lower, m) {
			return true
		}
	}
	return false
}

// looksLikeCode uses heuristics to decide if the comment text looks like
// commented-out Go code.
func looksLikeCode(text string) bool {
	// Try to parse the text as a Go statement/expression.
	// We wrap it in a function body to make the parser happy.
	src := "package _; func _() {\n" + text + "\n}"
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, 0)
	if err != nil {
		return false
	}

	// Successfully parsed: check if we got at least one statement in the
	// function body.
	for _, decl := range f.Decls {
		fn, ok := decl.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			continue
		}
		if len(fn.Body.List) > 0 {
			// Successfully parsed as valid Go statement(s).
			// Apply additional filters to reduce false positives.
			return !isLikelyDocumentation(text, fn.Body.List)
		}
	}

	return false
}

// isLikelyDocumentation applies additional heuristics on successfully parsed
// statements to determine if they are more likely documentation than code.
func isLikelyDocumentation(text string, stmts []ast.Stmt) bool {
	// If it parsed as a single expression statement that is just an
	// identifier, it is probably documentation (e.g. "// SomeType").
	if len(stmts) == 1 {
		if exprStmt, ok := stmts[0].(*ast.ExprStmt); ok {
			switch x := exprStmt.X.(type) {
			case *ast.Ident:
				// Single identifier like "SomeType" -- probably not code.
				return true
			case *ast.SelectorExpr:
				// Something like "pkg.Func" without parens is probably a reference.
				return true
			case *ast.BasicLit:
				// Literal alone ("123" or "hello") is ambiguous but lean towards doc.
				return true
			case *ast.BinaryExpr:
				// Simple binary like "a + b" might be documentation
				// Only flag if it has function calls or assignments
				_ = x
				return true
			}
		}

		// A single labeled statement could be a sentence ending with a colon.
		// e.g. "// Note: something" would parse as label "Note" with stmt.
		if _, ok := stmts[0].(*ast.LabeledStmt); ok {
			return true
		}
	}

	return false
}

// truncate shortens a string to at most maxLen characters, appending "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Name returns the rule name.
func (*NoCommentedOutCodeRule) Name() string {
	return "noCommentedOutCode"
}

// Group returns the rule group.
func (*NoCommentedOutCodeRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoCommentedOutCodeRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
