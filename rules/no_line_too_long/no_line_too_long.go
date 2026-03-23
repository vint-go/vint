package no_line_too_long

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strings"
	"unicode/utf8"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoLineTooLongRule reports lines that exceed a configured maximum character length.
type NoLineTooLongRule struct {
	lineLength int
	tabWidth   int
}

const (
	defaultLineLength = 120
	defaultTabWidth   = 1
)

// Configure validates and applies the rule configuration.
//
// Configuration implements the [lint.ConfigurableRule] interface.
func (r *NoLineTooLongRule) Configure(arguments lint.Arguments) error {
	r.lineLength = defaultLineLength
	r.tabWidth = defaultTabWidth

	if len(arguments) < 1 {
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noLineTooLong" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		switch {
		case isRuleOption(k, "line-length"):
			val, ok := lint.ToInt64(v)
			if !ok {
				return fmt.Errorf(`invalid configuration value for line-length in "noLineTooLong" rule; need integer but got %T`, v)
			}
			r.lineLength = int(val)
		case isRuleOption(k, "tab-width"):
			val, ok := lint.ToInt64(v)
			if !ok {
				return fmt.Errorf(`invalid configuration value for tab-width in "noLineTooLong" rule; need integer but got %T`, v)
			}
			r.tabWidth = int(val)
		}
	}

	return nil
}

// Apply applies the rule to given file.
func (r *NoLineTooLongRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	maxLen := r.lineLength
	if maxLen == 0 {
		maxLen = defaultLineLength
	}
	tabWidth := r.tabWidth
	if tabWidth == 0 {
		tabWidth = defaultTabWidth
	}

	// Compute import block line ranges to exclude.
	importRanges := importBlockRanges(file)

	var failures []lint.Failure
	spaces := strings.Repeat(" ", tabWidth)
	scanner := bufio.NewScanner(bytes.NewReader(file.Content()))
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip Go compiler directives (//go:...).
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//go:") {
			continue
		}

		// Skip lines inside import (...) blocks.
		if isInImportBlock(lineNum, importRanges) {
			continue
		}

		// Expand tabs and measure rune count.
		expanded := strings.ReplaceAll(line, "\t", spaces)
		runeCount := utf8.RuneCountInString(expanded)

		if runeCount > maxLen {
			failures = append(failures, lint.Failure{
				Category:   lint.FailureCategoryStyle,
				Confidence: 1,
				Position: lint.FailurePosition{
					Start: token.Position{
						Filename: file.Name,
						Line:     lineNum,
						Column:   1,
					},
					End: token.Position{
						Filename: file.Name,
						Line:     lineNum,
						Column:   runeCount,
					},
				},
				Failure: fmt.Sprintf("line is %d characters, exceeds limit of %d", runeCount, maxLen),
			})
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoLineTooLongRule) Name() string {
	return "noLineTooLong"
}

// Group returns the rule group.
func (*NoLineTooLongRule) Group() string {
	return "style"
}

// CacheTier returns the cache tier for this rule.
func (*NoLineTooLongRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// lineRange represents a range of line numbers (inclusive).
type lineRange struct {
	start int
	end   int
}

// importBlockRanges returns the line ranges of multi-line import blocks.
func importBlockRanges(file *lint.File) []lineRange {
	var ranges []lineRange
	for _, decl := range file.AST.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.IMPORT {
			continue
		}
		// Only multi-line import blocks (those with parentheses).
		if genDecl.Lparen == token.NoPos {
			continue
		}
		startPos := file.ToPosition(genDecl.Pos())
		endPos := file.ToPosition(genDecl.End())
		ranges = append(ranges, lineRange{start: startPos.Line, end: endPos.Line})
	}
	return ranges
}

// isInImportBlock returns true if the given line number falls within any import block range.
func isInImportBlock(line int, ranges []lineRange) bool {
	for _, r := range ranges {
		if line >= r.start && line <= r.end {
			return true
		}
	}
	return false
}

// isRuleOption returns true if arg and name are the same after normalization.
func isRuleOption(arg, name string) bool {
	return normalizeRuleOption(arg) == normalizeRuleOption(name)
}

// normalizeRuleOption returns an option name lowercased and without hyphens.
func normalizeRuleOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
