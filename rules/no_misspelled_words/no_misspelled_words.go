package no_misspelled_words

import (
	"fmt"
	"go/token"
	"strings"

	"github.com/golangci/misspell"
	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoMisspelledWordsRule detects commonly misspelled English words in Go source files.
type NoMisspelledWordsRule struct {
	locale      string
	mode        string
	extraWords  []extraWord
	ignoreRules []string
	replacer    *misspell.Replacer
}

type extraWord struct {
	Typo       string
	Correction string
}

// Configure validates and applies the rule configuration.
func (r *NoMisspelledWordsRule) Configure(arguments lint.Arguments) error {
	r.locale = ""
	r.mode = ""
	r.extraWords = nil
	r.ignoreRules = nil

	if len(arguments) < 1 {
		r.replacer = r.buildReplacer()
		return nil
	}

	argKV, ok := arguments[0].(map[string]any)
	if !ok {
		return fmt.Errorf(`invalid argument to the "noMisspelledWords" rule, expecting a k,v map, got %T`, arguments[0])
	}

	for k, v := range argKV {
		switch normalizeOption(k) {
		case "locale":
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf(`invalid configuration value for locale in "noMisspelledWords" rule; need string but got %T`, v)
			}
			s = strings.ToUpper(s)
			if s != "" && s != "US" && s != "UK" && s != "GB" {
				return fmt.Errorf(`invalid locale %q in "noMisspelledWords" rule; expected "US", "UK", "GB", or ""`, s)
			}
			r.locale = s
		case "mode":
			s, ok := v.(string)
			if !ok {
				return fmt.Errorf(`invalid configuration value for mode in "noMisspelledWords" rule; need string but got %T`, v)
			}
			r.mode = strings.ToLower(s)
		case "extrawords":
			list, ok := v.([]any)
			if !ok {
				return fmt.Errorf(`invalid configuration value for extra-words in "noMisspelledWords" rule; need list but got %T`, v)
			}
			for _, item := range list {
				m, ok := item.(map[string]any)
				if !ok {
					return fmt.Errorf(`invalid extra-words entry in "noMisspelledWords" rule; need map but got %T`, item)
				}
				typo, _ := m["typo"].(string)
				correction, _ := m["correction"].(string)
				if typo == "" || correction == "" {
					return fmt.Errorf(`extra-words entry must have non-empty "typo" and "correction" fields`)
				}
				r.extraWords = append(r.extraWords, extraWord{
					Typo:       strings.ToLower(typo),
					Correction: strings.ToLower(correction),
				})
			}
		case "ignorerules":
			list, ok := v.([]any)
			if !ok {
				return fmt.Errorf(`invalid configuration value for ignore-rules in "noMisspelledWords" rule; need list but got %T`, v)
			}
			for _, item := range list {
				s, ok := item.(string)
				if !ok {
					return fmt.Errorf(`invalid ignore-rules entry in "noMisspelledWords" rule; need string but got %T`, item)
				}
				r.ignoreRules = append(r.ignoreRules, s)
			}
		}
	}

	r.replacer = r.buildReplacer()
	return nil
}

// buildReplacer creates and compiles a misspell.Replacer from the current configuration.
func (r *NoMisspelledWordsRule) buildReplacer() *misspell.Replacer {
	replacer := &misspell.Replacer{
		Replacements: misspell.DictMain,
	}

	switch r.locale {
	case "US":
		replacer.AddRuleList(misspell.DictAmerican)
	case "UK", "GB":
		replacer.AddRuleList(misspell.DictBritish)
	}

	if len(r.extraWords) > 0 {
		additions := make([]string, 0, len(r.extraWords)*2)
		for _, ew := range r.extraWords {
			additions = append(additions, ew.Typo, ew.Correction)
		}
		replacer.AddRuleList(additions)
	}

	if len(r.ignoreRules) > 0 {
		replacer.RemoveRule(r.ignoreRules)
	}

	replacer.Compile()
	return replacer
}

// Apply applies the rule to given file.
func (r *NoMisspelledWordsRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	// In the default (Go-aware) mode, skip files with no comments.
	if r.mode != "unrestricted" && file.AST != nil && len(file.AST.Comments) == 0 {
		return nil
	}

	var failures []lint.Failure

	content := string(file.Content())

	var diffs []misspell.Diff
	if r.mode == "unrestricted" {
		_, diffs = r.replacer.Replace(content)
	} else { // default or "restricted" — both use Go-aware scanning
		_, diffs = r.replacer.ReplaceGo(content)
	}

	for _, diff := range diffs {
		failures = append(failures, lint.Failure{
			Category:   lint.FailureCategoryContent,
			Confidence: 1,
			Failure:    fmt.Sprintf(`"%s" is a misspelling of "%s"`, diff.Original, diff.Corrected),
			Position: lint.FailurePosition{
				Start: token.Position{
					Filename: file.Name,
					Line:     diff.Line,
					Column:   diff.Column,
				},
				End: token.Position{
					Filename: file.Name,
					Line:     diff.Line,
					Column:   diff.Column + len(diff.Original),
				},
			},
			ReplacementLine: diff.FullLine,
		})
	}

	return failures
}

// Name returns the rule name.
func (*NoMisspelledWordsRule) Name() string {
	return "noMisspelledWords"
}

// Group returns the rule group.
func (*NoMisspelledWordsRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoMisspelledWordsRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// normalizeOption returns an option name lowercased and without hyphens.
func normalizeOption(arg string) string {
	return strings.ToLower(strings.ReplaceAll(arg, "-", ""))
}
