package no_secret_in_serialization

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"

	"github.com/vint-go/vint/internal/rulecache"
	"github.com/vint-go/vint/lint"
)

// NoSecretInSerializationRule detects struct fields with secret-like names
// that are not excluded from JSON, YAML, XML, or TOML serialization output.
type NoSecretInSerializationRule struct{}

// Apply applies the rule to given file.
func (r *NoSecretInSerializationRule) Apply(file *lint.File, _ lint.Arguments) []lint.Failure {
	var failures []lint.Failure

	w := &lintSecretInSerialization{
		onFailure: func(f lint.Failure) {
			failures = append(failures, f)
		},
	}
	ast.Walk(w, file.AST)

	return failures
}

// Name returns the rule name.
func (*NoSecretInSerializationRule) Name() string {
	return "noSecretInSerialization"
}

// Group returns the rule group.
func (*NoSecretInSerializationRule) Group() string {
	return "security"
}

// CacheTier returns the cache tier for this rule.
func (*NoSecretInSerializationRule) CacheTier() rulecache.CacheTier {
	return rulecache.TierFileOnly
}

// secretFieldPattern matches field names that look like they contain sensitive data.
var secretFieldPattern = regexp.MustCompile(
	`(?i)(passwd|password|pwd|secret|token|api[_]?key|apikey|access[_]?key|auth[_]?token|private[_]?key|credentials?)`,
)

// serializationTags are the struct tag keys used for serialization.
var serializationTags = []string{"json", "yaml", "xml", "toml"}

type lintSecretInSerialization struct {
	onFailure func(lint.Failure)
}

func (w *lintSecretInSerialization) Visit(node ast.Node) ast.Visitor {
	structType, ok := node.(*ast.StructType)
	if !ok {
		return w
	}

	if structType.Fields == nil || structType.Fields.NumFields() == 0 {
		return w
	}

	for _, field := range structType.Fields.List {
		w.checkField(field)
	}

	return w
}

func (w *lintSecretInSerialization) checkField(field *ast.Field) {
	// Get the field name
	if len(field.Names) == 0 {
		return // embedded field, skip
	}

	fieldName := field.Names[0].Name
	if !secretFieldPattern.MatchString(fieldName) {
		return // not a secret-like field name
	}

	// Check if the field has serialization tags
	if field.Tag == nil {
		return // no tags at all — rule only flags fields with explicit serialization tags
	}

	tagValue := strings.Trim(field.Tag.Value, "`")
	if tagValue == "" {
		return
	}

	// Check each serialization tag
	for _, tagKey := range serializationTags {
		tagContent := getTagValue(tagValue, tagKey)
		if tagContent == "" {
			continue // tag not present
		}

		// Extract the name portion (before the first comma)
		name := tagContent
		if idx := strings.Index(tagContent, ","); idx != -1 {
			name = tagContent[:idx]
		}

		// If the name is "-", the field is excluded from serialization
		if name == "-" {
			continue
		}

		// Secret field is exposed via serialization
		w.onFailure(lint.Failure{
			Confidence: 1,
			Node:       field,
			Category:   lint.FailureCategoryBadPractice,
			Failure:    fmt.Sprintf("field %q with secret-like name is exposed via %s serialization, use `%s:\"-\"` to exclude it", fieldName, tagKey, tagKey),
		})
	}
}

// getTagValue extracts the value for a given key from a struct tag string.
// This is a simplified version of reflect.StructTag.Lookup.
func getTagValue(tag, key string) string {
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

		// Scan to colon to find key
		i = 0
		for i < len(tag) && tag[i] > ' ' && tag[i] != '"' && tag[i] != ':' {
			i++
		}
		if i == 0 || i+1 >= len(tag) || tag[i] != ':' || tag[i+1] != '"' {
			break
		}
		name := tag[:i]
		tag = tag[i+1:]

		// Scan quoted string to find value
		i = 1
		for i < len(tag) && tag[i] != '"' {
			if tag[i] == '\\' {
				i++
			}
			i++
		}
		if i >= len(tag) {
			break
		}
		qvalue := tag[:i+1]
		tag = tag[i+1:]

		if key == name {
			// Unquote
			value := qvalue[1 : len(qvalue)-1]
			return value
		}
	}
	return ""
}
