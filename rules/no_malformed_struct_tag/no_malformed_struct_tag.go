package no_malformed_struct_tag

import (
	"fmt"
	"go/ast"
	"reflect"
	"strconv"
	"strings"

	"github.com/strowk/vint/internal/rulecache"
	"github.com/strowk/vint/lint"
)

// NoMalformedStructTagRule checks that struct field tags are well formed.
// It validates that struct field tags follow the correct syntax as defined
// by reflect.StructTag, detects duplicate tag keys, misuse of tag options,
// and incorrectly structured tag key-value pairs.
type NoMalformedStructTagRule struct{}

// Apply applies the rule to given file.
func (r *NoMalformedStructTagRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintMalformedStructTag{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoMalformedStructTagRule) Name() string {
	return "noMalformedStructTag"
}

// Group returns the rule group.
func (*NoMalformedStructTagRule) Group() string {
	return "correctness"
}

// CacheTier returns the cache tier for this rule.
func (*NoMalformedStructTagRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

type lintMalformedStructTag struct {
	onFailure func(lint.Failure)
}

func (w *lintMalformedStructTag) Visit(node ast.Node) ast.Visitor {
	st, ok := node.(*ast.StructType)
	if !ok {
		return w
	}

	if st.Fields == nil || st.Fields.NumFields() < 1 {
		return w
	}

	for _, field := range st.Fields.List {
		if field.Tag == nil {
			continue
		}
		w.checkFieldTag(field)
	}

	return w
}

// checkFieldTag validates a struct field's tag.
func (w *lintMalformedStructTag) checkFieldTag(field *ast.Field) {
	tagLit := field.Tag
	tagValue, err := strconv.Unquote(tagLit.Value)
	if err != nil {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryBadPractice,
			Confidence: 1,
			Node:       tagLit,
			Failure:    "malformed struct tag: cannot unquote tag value",
		})
		return
	}

	// Check for space before value (e.g., `json: "name"`)
	if msg := checkSpaceInTagValue(tagValue); msg != "" {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryBadPractice,
			Confidence: 1,
			Node:       tagLit,
			Failure:    msg,
		})
		return
	}

	// Use reflect.StructTag to validate basic syntax
	tag := reflect.StructTag(tagValue)

	// Try to parse the tag manually to detect syntax issues
	if msg := validateTagSyntax(tagValue); msg != "" {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryBadPractice,
			Confidence: 1,
			Node:       tagLit,
			Failure:    msg,
		})
		return
	}

	// Check for duplicate keys
	if msg := checkDuplicateKeys(tag, tagValue); msg != "" {
		w.onFailure(lint.Failure{
			Category:   lint.FailureCategoryBadPractice,
			Confidence: 1,
			Node:       tagLit,
			Failure:    msg,
		})
	}
}

// checkSpaceInTagValue checks for the pattern `key: "value"` which is a common mistake.
func checkSpaceInTagValue(tag string) string {
	rest := tag
	for rest != "" {
		// Skip leading spaces
		i := 0
		for i < len(rest) && rest[i] == ' ' {
			i++
		}
		rest = rest[i:]
		if rest == "" {
			break
		}

		// Scan key
		i = 0
		for i < len(rest) && rest[i] > ' ' && rest[i] != '"' && rest[i] != ':' {
			i++
		}
		if i == 0 {
			break
		}
		key := rest[:i]
		rest = rest[i:]

		// Check for colon followed by space and then a quote
		if strings.HasPrefix(rest, ": \"") {
			return fmt.Sprintf("malformed struct tag: key %q has space before value in tag", key)
		}

		// Skip to end of this key-value pair
		if len(rest) == 0 || rest[0] != ':' {
			break
		}
		rest = rest[1:] // skip colon

		if len(rest) == 0 || rest[0] != '"' {
			break
		}
		rest = rest[1:] // skip opening quote

		// Find closing quote
		for i = 0; i < len(rest); i++ {
			if rest[i] == '\\' {
				i++ // skip escaped char
				continue
			}
			if rest[i] == '"' {
				break
			}
		}
		if i >= len(rest) {
			break
		}
		rest = rest[i+1:]
	}
	return ""
}

// validateTagSyntax validates the basic syntax of a struct tag.
// It follows the same parsing rules as reflect.StructTag.
func validateTagSyntax(tag string) string {
	for tag != "" {
		// Skip leading spaces
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan key: non-space, non-colon, non-quote characters
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != '"' && tag[i] != ':' {
			i++
		}
		if i == 0 {
			return fmt.Sprintf("malformed struct tag: invalid character %q in tag key", tag[0])
		}
		if i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			return "malformed struct tag: missing colon or quote in tag"
		}

		// Skip past the key and colon
		tag = tag[i+1:]

		// Unquote the value
		val, rest, err := unquoteTag(tag)
		if err != nil {
			return "malformed struct tag: badly quoted value in tag"
		}
		_ = val
		tag = rest
	}
	return ""
}

// unquoteTag parses a quoted string at the start of tag, returning the
// unquoted value and the remainder of the tag string.
func unquoteTag(tag string) (string, string, error) {
	if len(tag) == 0 || tag[0] != '"' {
		return "", "", fmt.Errorf("expected opening quote")
	}

	// Find the closing quote
	for i := 1; i < len(tag); i++ {
		if tag[i] == '\\' {
			i++ // skip next character
			continue
		}
		if tag[i] == '"' {
			val, err := strconv.Unquote(tag[:i+1])
			if err != nil {
				return "", "", err
			}
			return val, tag[i+1:], nil
		}
	}

	return "", "", fmt.Errorf("missing closing quote")
}

// checkDuplicateKeys checks if the tag has duplicate keys.
func checkDuplicateKeys(tag reflect.StructTag, raw string) string {
	keys := map[string]bool{}
	rest := raw
	for rest != "" {
		// Skip spaces
		i := 0
		for i < len(rest) && rest[i] == ' ' {
			i++
		}
		rest = rest[i:]
		if rest == "" {
			break
		}

		// Scan key
		i = 0
		for i < len(rest) && rest[i] > ' ' && rest[i] != '"' && rest[i] != ':' {
			i++
		}
		if i == 0 {
			break
		}
		key := rest[:i]
		rest = rest[i:]

		if keys[key] {
			return fmt.Sprintf("malformed struct tag: duplicate tag key %q", key)
		}
		keys[key] = true

		// Skip past colon and quoted value
		if len(rest) == 0 || rest[0] != ':' {
			break
		}
		rest = rest[1:]

		if len(rest) == 0 || rest[0] != '"' {
			break
		}
		rest = rest[1:]

		// Find closing quote
		for i = 0; i < len(rest); i++ {
			if rest[i] == '\\' {
				i++
				continue
			}
			if rest[i] == '"' {
				break
			}
		}
		if i >= len(rest) {
			break
		}
		rest = rest[i+1:]
	}

	return ""
}
