// Package mcpserver implements an MCP (Model Context Protocol) server for vint.
// It exposes linting functionality as MCP tools over stdio transport.
package mcpserver

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/vint-go/vint/config"
	"github.com/vint-go/vint/lint"
	"github.com/vint-go/vint/vintlint0"
)

const (
	maxTopRules          = 5
	maxTopFiles          = 3
	maxDetailedRules     = 3
	maxFailuresPerRule   = 5
)

// reportArgs defines the input schema for the vint_report tool.
type reportArgs struct {
	// Package glob patterns to lint (e.g. "./cmd/...", "./internal/api"). Defaults to ".".
	Packages []string `json:"packages,omitempty" jsonschema:"Package glob patterns to lint. Defaults to current directory."`
	// File path substrings to filter results by (e.g. "handler.go", "api/").
	Files []string `json:"files,omitempty" jsonschema:"File path substrings to filter results by."`
	// Rule names to filter results by (e.g. "errcheck", "unused").
	Rules []string `json:"rules,omitempty" jsonschema:"Rule names to filter by. Partial match supported."`
	// Path to vint config file (TOML or YAML).
	ConfigPath string `json:"config_path,omitempty" jsonschema:"Path to vint configuration file."`
}

// Run starts the MCP server on stdio transport and blocks until the session ends.
func Run(ctx context.Context) error {
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "vint",
		Title:   "Vint Linter",
		Version: "1.0.0",
	}, &mcp.ServerOptions{
		Instructions: "Vint is a fast Go linter. Use the vint_report tool to get a linting report for Go packages.",
	})

	mcp.AddTool(server, &mcp.Tool{
		Name:        "vint_report",
		Description: "Run vint linter and produce a structured markdown report of lint failures. Supports filtering by package, file, and rule name. Output is capped for readability.",
	}, handleReport)

	return server.Run(ctx, &mcp.StdioTransport{})
}

func handleReport(_ context.Context, _ *mcp.CallToolRequest, args reportArgs) (*mcp.CallToolResult, any, error) {
	failures, ruleCount, err := runLint(args)
	if err != nil {
		return nil, nil, fmt.Errorf("linting failed: %w", err)
	}

	report := formatReport(failures, ruleCount, args)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: report},
		},
	}, nil, nil
}

// runLint executes the vint linter and collects filtered failures.
func runLint(args reportArgs) ([]lint.Failure, int, error) {
	conf, err := config.GetConfig(args.ConfigPath)
	if err != nil {
		return nil, 0, fmt.Errorf("loading config: %w", err)
	}

	lintingRules, err := config.GetLintingRules(conf, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("getting rules: %w", err)
	}

	includes := args.Packages
	if len(includes) == 0 {
		includes = []string{"."}
	}

	// Disable cache for MCP to always get fresh results.
	os.Setenv("VINT_NO_CACHE", "1")

	linter := vintlint0.New(lintingRules, *conf)
	failChan, err := linter.Lint(includes, []string{"vendor/..."})
	if err != nil {
		return nil, 0, err
	}

	ruleFilter := buildRuleFilter(args.Rules)
	fileFilter := args.Files

	var failures []lint.Failure
	for f := range failChan {
		if f.IsInternal() {
			continue
		}
		if f.Confidence < conf.Confidence {
			continue
		}
		if !matchesRuleFilter(f.RuleName, ruleFilter) {
			continue
		}
		if !matchesFileFilter(f.Filename(), fileFilter) {
			continue
		}
		failures = append(failures, f)
	}

	return failures, len(lintingRules), nil
}

func buildRuleFilter(rules []string) []string {
	if len(rules) == 0 {
		return nil
	}
	lower := make([]string, len(rules))
	for i, r := range rules {
		lower[i] = strings.ToLower(r)
	}
	return lower
}

func matchesRuleFilter(ruleName string, filter []string) bool {
	if len(filter) == 0 {
		return true
	}
	lower := strings.ToLower(ruleName)
	for _, f := range filter {
		if strings.Contains(lower, f) {
			return true
		}
	}
	return false
}

func matchesFileFilter(filename string, filter []string) bool {
	if len(filter) == 0 {
		return true
	}
	// Normalize to forward slashes for consistent matching.
	norm := filepath.ToSlash(filename)
	for _, f := range filter {
		if strings.Contains(norm, f) {
			return true
		}
	}
	return false
}

// formatReport produces a capped markdown report from the collected failures.
func formatReport(failures []lint.Failure, ruleCount int, args reportArgs) string {
	if len(failures) == 0 {
		return formatNoFailures(ruleCount, args)
	}

	var b strings.Builder

	// Aggregate stats.
	ruleFailCounts := map[string]int{}
	fileFailCounts := map[string]int{}
	for _, f := range failures {
		ruleFailCounts[f.RuleName]++
		fileFailCounts[f.Filename()]++
	}

	// Sort rules by count descending.
	type kv struct {
		key   string
		count int
	}
	topRules := sortedTopN(ruleFailCounts, maxTopRules)
	topFiles := sortedTopN(fileFailCounts, maxTopFiles)

	// --- Header / Summary ---
	b.WriteString("# Vint Lint Report\n\n")
	b.WriteString("## Summary\n\n")
	fmt.Fprintf(&b, "- **Total failures**: %d\n", len(failures))
	fmt.Fprintf(&b, "- **Files with failures**: %d\n", len(fileFailCounts))
	fmt.Fprintf(&b, "- **Failing rules**: %d\n", len(ruleFailCounts))
	fmt.Fprintf(&b, "- **Rules evaluated**: %d\n", ruleCount)
	if len(args.Packages) > 0 {
		fmt.Fprintf(&b, "- **Scope (packages)**: `%s`\n", strings.Join(args.Packages, "`, `"))
	}
	if len(args.Files) > 0 {
		fmt.Fprintf(&b, "- **Filter (files)**: `%s`\n", strings.Join(args.Files, "`, `"))
	}
	if len(args.Rules) > 0 {
		fmt.Fprintf(&b, "- **Filter (rules)**: `%s`\n", strings.Join(args.Rules, "`, `"))
	}
	b.WriteString("\n")

	// --- Top Failing Rules ---
	b.WriteString("## Top Failing Rules\n\n")
	b.WriteString("| # | Rule | Failures |\n")
	b.WriteString("|---|------|----------|\n")
	for i, kv := range topRules {
		fmt.Fprintf(&b, "| %d | `%s` | %d |\n", i+1, kv.key, kv.count)
	}
	if remaining := len(ruleFailCounts) - len(topRules); remaining > 0 {
		fmt.Fprintf(&b, "\n*...and %d more failing rule(s) not shown.*\n", remaining)
	}
	b.WriteString("\n")

	// --- Top Files with Failures ---
	b.WriteString("## Top Files with Failures\n\n")
	b.WriteString("| # | File | Failures |\n")
	b.WriteString("|---|------|----------|\n")
	for i, kv := range topFiles {
		fmt.Fprintf(&b, "| %d | `%s` | %d |\n", i+1, kv.key, kv.count)
	}
	if remaining := len(fileFailCounts) - len(topFiles); remaining > 0 {
		fmt.Fprintf(&b, "\n*...and %d more file(s) with failures not shown.*\n", remaining)
	}
	b.WriteString("\n")

	// --- Detailed Failures ---
	// Show detailed failures for up to maxDetailedRules rules,
	// preferring the top rules from the summary above.
	b.WriteString("## Failure Details\n\n")
	detailRules := topRules
	if len(detailRules) > maxDetailedRules {
		detailRules = detailRules[:maxDetailedRules]
	}

	for _, rk := range detailRules {
		fmt.Fprintf(&b, "### `%s` (%d total failure(s))\n\n", rk.key, rk.count)

		shown := 0
		for _, f := range failures {
			if f.RuleName != rk.key {
				continue
			}
			if shown >= maxFailuresPerRule {
				break
			}
			shown++

			pos := f.Position.Start
			fmt.Fprintf(&b, "**%s:%d:%d**", filepath.ToSlash(pos.Filename), pos.Line, pos.Column)
			if f.Position.End.Line > 0 && f.Position.End.Line != pos.Line {
				fmt.Fprintf(&b, " → :%d:%d", f.Position.End.Line, f.Position.End.Column)
			}
			b.WriteString("\n")
			fmt.Fprintf(&b, "> %s\n", f.Failure)
			if f.Category != "" {
				fmt.Fprintf(&b, "> Category: `%s`\n", f.Category)
			}
			if f.ReplacementLine != "" {
				fmt.Fprintf(&b, "> Suggested fix: `%s`\n", f.ReplacementLine)
			}
			b.WriteString("\n")
		}

		if rk.count > maxFailuresPerRule {
			fmt.Fprintf(&b, "*...%d more failure(s) for this rule not shown.*\n\n", rk.count-maxFailuresPerRule)
		}
	}

	if len(ruleFailCounts) > maxDetailedRules {
		fmt.Fprintf(&b, "*Details for %d more rule(s) omitted. Use the `rules` filter to see specific rules.*\n", len(ruleFailCounts)-maxDetailedRules)
	}

	return b.String()
}

func formatNoFailures(ruleCount int, args reportArgs) string {
	var b strings.Builder
	b.WriteString("# Vint Lint Report\n\n")
	b.WriteString("## Summary\n\n")
	b.WriteString("**No lint failures found.**\n\n")
	fmt.Fprintf(&b, "- **Rules evaluated**: %d\n", ruleCount)
	if len(args.Packages) > 0 {
		fmt.Fprintf(&b, "- **Scope (packages)**: `%s`\n", strings.Join(args.Packages, "`, `"))
	}
	if len(args.Files) > 0 {
		fmt.Fprintf(&b, "- **Filter (files)**: `%s`\n", strings.Join(args.Files, "`, `"))
	}
	if len(args.Rules) > 0 {
		fmt.Fprintf(&b, "- **Filter (rules)**: `%s`\n", strings.Join(args.Rules, "`, `"))
	}
	return b.String()
}

type kv struct {
	key   string
	count int
}

func sortedTopN(m map[string]int, n int) []kv {
	items := make([]kv, 0, len(m))
	for k, v := range m {
		items = append(items, kv{k, v})
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].count != items[j].count {
			return items[i].count > items[j].count
		}
		return items[i].key < items[j].key
	})
	if len(items) > n {
		items = items[:n]
	}
	return items
}
