package vintlint0

import (
	"fmt"
	"go/token"
	"math"
	"regexp"
	"strings"

	"github.com/vint-go/vint/lint"
)

var vintIgnoreRegexp = regexp.MustCompile(`^//\s*vint-ignore(-all|-start|-end)?\s+(\S+?)(?::\s*(.+))?$`)

const (
	vintIgnoreModifierGroup    = 1
	vintIgnoreRuleNameGroup    = 2
	vintIgnoreExplanationGroup = 3
)

type pendingStart struct {
	line int
}

// parseVintIgnoreSuppressions scans comments in f for vint-ignore directives
// and returns a map of disabled intervals keyed by rule name.
//
// Supported directives:
//
//	// vint-ignore <rule>: <explanation>        — suppress next line
//	// vint-ignore-all <rule>: <explanation>    — suppress entire file (must be before package decl)
//	// vint-ignore-start <rule>: <explanation>  — begin suppression range
//	// vint-ignore-end <rule>                   — end suppression range
func parseVintIgnoreSuppressions(f *lint.File, failures chan<- lint.Failure) map[string][]lint.DisabledInterval {
	result := map[string][]lint.DisabledInterval{}
	// Track pending -start directives: ruleName -> stack of start lines.
	pending := map[string][]pendingStart{}

	for _, cg := range f.AST.Comments {
		for _, c := range cg.List {
			match := vintIgnoreRegexp.FindStringSubmatch(c.Text)
			if len(match) == 0 {
				continue
			}

			modifier := match[vintIgnoreModifierGroup]
			ruleName := match[vintIgnoreRuleNameGroup]
			explanation := strings.TrimSpace(match[vintIgnoreExplanationGroup])
			commentLine := f.ToPosition(c.Pos()).Line

			switch modifier {
			case "": // inline: suppress next line
				if explanation == "" {
					failures <- lint.Failure{
						Confidence: 1,
						RuleName:   "vint-ignore",
						Failure:    "vint-ignore directive requires an explanation",
						Position:   lint.ToFailurePosition(c.Pos(), c.End(), f),
					}
					continue
				}
				nextLine := commentLine + 1
				result[ruleName] = append(result[ruleName], lint.DisabledInterval{
					RuleName: ruleName,
					From:     posAt(f.Name, nextLine),
					To:       posAt(f.Name, nextLine),
				})

			case "-all": // file-level suppression
				if explanation == "" {
					failures <- lint.Failure{
						Confidence: 1,
						RuleName:   "vint-ignore",
						Failure:    "vint-ignore-all directive requires an explanation",
						Position:   lint.ToFailurePosition(c.Pos(), c.End(), f),
					}
					continue
				}
				if c.Pos() >= f.AST.Package {
					failures <- lint.Failure{
						Confidence: 1,
						RuleName:   "vint-ignore",
						Failure:    "vint-ignore-all must appear before the package declaration",
						Position:   lint.ToFailurePosition(c.Pos(), c.End(), f),
					}
					continue
				}
				result[ruleName] = append(result[ruleName], lint.DisabledInterval{
					RuleName: ruleName,
					From:     posAt(f.Name, 1),
					To:       posAt(f.Name, math.MaxInt32),
				})

			case "-start":
				if explanation == "" {
					failures <- lint.Failure{
						Confidence: 1,
						RuleName:   "vint-ignore",
						Failure:    "vint-ignore-start directive requires an explanation",
						Position:   lint.ToFailurePosition(c.Pos(), c.End(), f),
					}
					continue
				}
				pending[ruleName] = append(pending[ruleName], pendingStart{line: commentLine})

			case "-end":
				stack := pending[ruleName]
				if len(stack) == 0 {
					failures <- lint.Failure{
						Confidence: 1,
						RuleName:   "vint-ignore",
						Failure:    fmt.Sprintf("vint-ignore-end for %s without matching vint-ignore-start", ruleName),
						Position:   lint.ToFailurePosition(c.Pos(), c.End(), f),
					}
					continue
				}
				start := stack[len(stack)-1]
				pending[ruleName] = stack[:len(stack)-1]
				result[ruleName] = append(result[ruleName], lint.DisabledInterval{
					RuleName: ruleName,
					From:     posAt(f.Name, start.line),
					To:       posAt(f.Name, commentLine),
				})
			}
		}
	}

	// Emit failures for unmatched -start directives.
	for ruleName, stack := range pending {
		for _, s := range stack {
			failures <- lint.Failure{
				Confidence: 1,
				RuleName:   "vint-ignore",
				Failure:    fmt.Sprintf("vint-ignore-start for %s at line %d has no matching vint-ignore-end", ruleName, s.line),
				Position: lint.FailurePosition{
					Start: posAt(f.Name, s.line),
					End:   posAt(f.Name, s.line),
				},
			}
		}
	}

	return result
}

func posAt(filename string, line int) token.Position {
	return token.Position{Filename: filename, Line: line}
}
