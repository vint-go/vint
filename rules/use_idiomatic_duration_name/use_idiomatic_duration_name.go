package use_idiomatic_duration_name

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// UseIdiomaticDurationNameRule checks that variables of type time.Duration
// do not have names that imply a specific time unit (e.g. timeoutSecs, delayMs).
type UseIdiomaticDurationNameRule struct{}

// Apply applies the rule to given file.
func (r *UseIdiomaticDurationNameRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintDurationName{file: file, onFailure: onFailure}
	ast.Walk(w, file.AST)

	return failures
}

// ApplyToNode applies the rule while walking the AST together with other rules.
func (r *UseIdiomaticDurationNameRule) ApplyToNode(file *lint.File, node ast.Node, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	file.Pkg.TypeCheck()

	onFailure := func(f lint.Failure) {
		failures = append(failures, f)
	}

	w := &lintDurationName{file: file, onFailure: onFailure}
	w.Visit(node)

	return failures
}

// Name returns the rule name.
func (*UseIdiomaticDurationNameRule) Name() string {
	return "useIdiomaticDurationName"
}

// Group returns the rule group.
func (*UseIdiomaticDurationNameRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*UseIdiomaticDurationNameRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierPackageAware
}

// RequiresTypecheck returns true because this rule uses type information.
func (*UseIdiomaticDurationNameRule) RequiresTypecheck() bool {
	return true
}

// timeSuffixes is a list of name suffixes that imply a time unit.
var timeSuffixes = []string{
	"Hour", "Hours",
	"Min", "Mins", "Minutes", "Minute",
	"Sec", "Secs", "Seconds", "Second",
	"Msec", "Msecs",
	"Milli", "Millis", "Milliseconds", "Millisecond",
	"Usec", "Usecs", "Microseconds", "Microsecond",
	"MS", "Ms",
}

type lintDurationName struct {
	file      *lint.File
	onFailure func(lint.Failure)
}

func (w *lintDurationName) Visit(node ast.Node) ast.Visitor {
	v, ok := node.(*ast.ValueSpec)
	if !ok {
		return w
	}

	for _, name := range v.Names {
		origTyp := w.file.Pkg.TypeOf(name)
		if origTyp == nil {
			continue
		}

		// Look for time.Duration or *time.Duration;
		// the latter is common when using flag.Duration.
		typ := origTyp
		if pt, ok := typ.(*types.Pointer); ok {
			typ = pt.Elem()
		}

		if !isNamedType(typ, "time", "Duration") {
			continue
		}

		suffix := ""
		for _, suf := range timeSuffixes {
			if strings.HasSuffix(name.Name, suf) {
				suffix = suf
				break
			}
		}
		if suffix == "" {
			continue
		}

		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryNaming,
			Confidence: 0.9,
			Node:       v,
			Failure:    fmt.Sprintf("var %s is of type %v; don't use unit-specific suffix %q", name.Name, origTyp, suffix),
		})
	}

	return w
}

func isNamedType(typ types.Type, importPath, name string) bool {
	n, ok := typ.(*types.Named)
	if !ok {
		return false
	}

	typeName := n.Obj()
	return typeName != nil && typeName.Pkg() != nil && typeName.Pkg().Path() == importPath && typeName.Name() == name
}
