package no_trojan_source_bidi

import (
	"bufio"
	"bytes"
	"fmt"
	"go/token"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// bidiChars maps Unicode bidirectional control characters to their names.
var bidiChars = map[rune]string{
	'\u200E': "Left-to-Right Mark",
	'\u200F': "Right-to-Left Mark",
	'\u202A': "Left-to-Right Embedding",
	'\u202B': "Right-to-Left Embedding",
	'\u202C': "Pop Directional Formatting",
	'\u202D': "Left-to-Right Override",
	'\u202E': "Right-to-Left Override",
	'\u2066': "Left-to-Right Isolate",
	'\u2067': "Right-to-Left Isolate",
	'\u2068': "First Strong Isolate",
	'\u2069': "Pop Directional Isolate",
}

// NoTrojanSourceBidiRule detects Trojan Source attacks using bidirectional
// Unicode control characters in Go source files.
type NoTrojanSourceBidiRule struct{}

// Apply applies the rule to given file.
func (r *NoTrojanSourceBidiRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	content := file.Content()
	scanner := bufio.NewScanner(bytes.NewReader(content))
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		for _, ch := range line {
			name, isBidi := bidiChars[ch]
			if !isBidi {
				continue
			}
			col := strings.IndexRune(line, ch)
			failures = append(failures, lint.Failure{
				Category:   lint.FailureCategoryContent,
				Confidence: 1,
				Failure:    fmt.Sprintf("found bidirectional Unicode control character %s (U+%04X)", name, ch),
				Position: lint.FailurePosition{
					Start: token.Position{
						Filename: file.Name,
						Line:     lineNum,
						Column:   col + 1,
					},
					End: token.Position{
						Filename: file.Name,
						Line:     lineNum,
						Column:   col + 2,
					},
				},
			})
			break // report only the first BiDi character per line
		}
	}

	return failures
}

// Name returns the rule name.
func (*NoTrojanSourceBidiRule) Name() string {
	return "noTrojanSourceBidi"
}

// Group returns the rule group.
func (*NoTrojanSourceBidiRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoTrojanSourceBidiRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}
