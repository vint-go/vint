package use_errorf

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"

	"github.com/vint-go/vint/internal/astutils"
	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// ErrorfRule suggests using `fmt.Errorf` instead of `errors.New(fmt.Sprintf())`.
type ErrorfRule struct{}

// Apply applies the rule to given file.
func (*ErrorfRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	fileAst := file.AST
	walker := lintErrorf{
		file:    file,
		fileAst: fileAst,
		onFailure: func(failure lint.Failure) {
			failures = append(failures, failure)
		},
	}

	file.Pkg.TypeCheck()
	ast.Walk(walker, fileAst)

	return failures
}

// Name returns the rule name.
func (*ErrorfRule) Name() string {
	return "useErrorf"
}

// Group returns the group the rule belongs to.
func (*ErrorfRule) Group() string {
	return "style"
}

// RequiresTypecheck returns true because this rule uses type information.
func (*ErrorfRule) RequiresTypecheck() bool { return true }

type lintErrorf struct {
	file      *lint.File
	fileAst   *ast.File
	onFailure func(lint.Failure)
}

func (w lintErrorf) Visit(n ast.Node) ast.Visitor {
	ce, ok := n.(*ast.CallExpr)
	if !ok || len(ce.Args) != 1 {
		return w
	}
	isErrorsNew := astutils.IsPkgDotName(ce.Fun, "errors", "New")
	var isTestingError bool
	se, ok := ce.Fun.(*ast.SelectorExpr)
	if ok && se.Sel.Name == "Error" {
		if typ := w.file.Pkg.TypeOf(se.X); typ != nil {
			isTestingError = typ.String() == "*testing.T"
		}
	}
	if !isErrorsNew && !isTestingError {
		return w
	}
	arg := ce.Args[0]
	ce, ok = arg.(*ast.CallExpr)
	if !ok || !astutils.IsPkgDotName(ce.Fun, "fmt", "Sprintf") {
		return w
	}
	errorfPrefix := "fmt"
	if isTestingError {
		errorfPrefix = w.file.Render(se.X)
	}

	failure := lint.Failure{
		Category:   lint.FailureCategoryErrors,
		Node:       n,
		Confidence: 1,
		Failure:    fmt.Sprintf("should replace %s(fmt.Sprintf(...)) with %s.Errorf(...)", w.file.Render(se), errorfPrefix),
	}

	m := srcLineWithMatch(w.file, ce, `^(.*)`+w.file.Render(se)+`\(fmt\.Sprintf\((.*)\)\)(.*)$`)
	if m != nil {
		failure.ReplacementLine = m[1] + errorfPrefix + ".Errorf(" + m[2] + ")" + m[3]
	}

	w.onFailure(failure)

	return w
}

func srcLineWithMatch(file *lint.File, node ast.Node, pattern string) (m []string) {
	line := srcLine(file.Content(), file.ToPosition(node.Pos()))
	line = strings.TrimSuffix(line, "\n")
	rx := regexp.MustCompile(pattern)
	return rx.FindStringSubmatch(line)
}

// CacheTier returns the cache tier for this rule.
func (*ErrorfRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

func srcLine(src []byte, p token.Position) string {
	// Run to end of line in both directions if not at line start/end.
	lo, hi := p.Offset, p.Offset+1
	for lo > 0 && src[lo-1] != '\n' {
		lo--
	}
	for hi < len(src) && src[hi-1] != '\n' {
		hi++
	}
	return string(src[lo:hi])
}
