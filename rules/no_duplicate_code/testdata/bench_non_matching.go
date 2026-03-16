package fixtures

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"
	"time"
)

// This file contains 15 large, structurally-different functions.
// Each has 160+ tokens but different AST structure, so no pair
// should match the 150-token duplicate threshold. This forces
// the DP to run to full completion for every pair.

func nmHandleHTTPRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.URL.Path == "/health" {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
		return
	}
	if r.Header.Get("Authorization") == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		http.Error(w, "bad content type", http.StatusBadRequest)
		return
	}
	query := r.URL.Query()
	page := query.Get("page")
	if page == "" {
		page = "1"
	}
	limit := query.Get("limit")
	if limit == "" {
		limit = "10"
	}
	accept := r.Header.Get("Accept")
	if accept == "" {
		accept = "application/json"
	}
	userAgent := r.Header.Get("User-Agent")
	if strings.Contains(userAgent, "bot") {
		http.Error(w, "bots not allowed", http.StatusForbidden)
		return
	}
	_, _ = fmt.Fprintf(w, "page=%s&limit=%s&accept=%s", page, limit, accept)
}

func nmProcessCSVData(lines []string) ([][]string, error) {
	var result [][]string
	for i, line := range lines {
		if i == 0 {
			continue // skip header
		}
		fields := strings.Split(line, ",")
		if len(fields) < 3 {
			return nil, fmt.Errorf("line %d: expected 3 fields, got %d", i, len(fields))
		}
		trimmed := make([]string, len(fields))
		for j, f := range fields {
			trimmed[j] = strings.TrimSpace(f)
		}
		if trimmed[0] == "" {
			return nil, fmt.Errorf("line %d: empty first field", i)
		}
		if strings.Contains(trimmed[1], "\"") {
			trimmed[1] = strings.ReplaceAll(trimmed[1], "\"", "")
		}
		if strings.HasPrefix(trimmed[2], " ") {
			trimmed[2] = strings.TrimLeft(trimmed[2], " ")
		}
		if len(trimmed[0]) > 100 {
			return nil, fmt.Errorf("line %d: first field too long", i)
		}
		if len(trimmed[1]) > 200 {
			return nil, fmt.Errorf("line %d: second field too long", i)
		}
		result = append(result, trimmed)
	}
	if len(result) == 0 {
		return nil, errors.New("no data rows found")
	}
	return result, nil
}

func nmSortFilterAndDeduplicate(items []string, minLen int, prefix string) []string {
	if len(items) == 0 {
		return nil
	}
	filtered := make([]string, 0, len(items))
	for _, item := range items {
		if len(item) < minLen {
			continue
		}
		if prefix != "" && !strings.HasPrefix(item, prefix) {
			continue
		}
		lower := strings.ToLower(item)
		if strings.Contains(lower, "deprecated") {
			continue
		}
		if strings.ContainsAny(lower, "!@#$%") {
			continue
		}
		filtered = append(filtered, item)
	}
	sort.Strings(filtered)
	result := make([]string, 0, len(filtered))
	if len(filtered) > 0 {
		result = append(result, filtered[0])
	}
	for i := 1; i < len(filtered); i++ {
		if filtered[i] != filtered[i-1] {
			result = append(result, filtered[i])
		}
	}
	for i, s := range result {
		result[i] = strings.TrimSpace(s)
	}
	return result
}

func nmBuildHTMLTable(headers []string, rows [][]string, classes map[string]string) string {
	var b strings.Builder
	b.WriteString("<table class=\"data-table\">\n<thead>\n<tr>\n")
	for _, h := range headers {
		cls := classes[h]
		if cls != "" {
			b.WriteString("<th class=\"")
			b.WriteString(cls)
			b.WriteString("\">")
		} else {
			b.WriteString("<th>")
		}
		b.WriteString(h)
		b.WriteString("</th>\n")
	}
	b.WriteString("</tr>\n</thead>\n<tbody>\n")
	for rowIdx, row := range rows {
		if rowIdx%2 == 0 {
			b.WriteString("<tr class=\"even\">\n")
		} else {
			b.WriteString("<tr class=\"odd\">\n")
		}
		for _, cell := range row {
			b.WriteString("<td>")
			escaped := strings.ReplaceAll(cell, "&", "&amp;")
			escaped = strings.ReplaceAll(escaped, "<", "&lt;")
			escaped = strings.ReplaceAll(escaped, ">", "&gt;")
			b.WriteString(escaped)
			b.WriteString("</td>\n")
		}
		b.WriteString("</tr>\n")
	}
	b.WriteString("</tbody>\n</table>")
	return b.String()
}

func nmRetryWithBackoff(fn func() error, maxRetries int, initialDelay time.Duration) error {
	var lastErr error
	delay := initialDelay
	if delay == 0 {
		delay = time.Millisecond * 100
	}
	maxDelay := time.Second * 30
	for attempt := 0; attempt < maxRetries; attempt++ {
		lastErr = fn()
		if lastErr == nil {
			return nil
		}
		errStr := lastErr.Error()
		if strings.Contains(errStr, "permanent") {
			return fmt.Errorf("permanent failure on attempt %d: %w", attempt, lastErr)
		}
		if strings.Contains(errStr, "timeout") {
			delay = delay / 2
		}
		if attempt < maxRetries-1 {
			time.Sleep(delay)
			delay *= 2
			if delay > maxDelay {
				delay = maxDelay
			}
		}
		_ = fmt.Sprintf("attempt %d/%d failed: %v", attempt+1, maxRetries, lastErr)
	}
	return fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

func nmConcurrentFetch(urls []string, timeout time.Duration) map[string]string {
	var mu sync.Mutex
	results := make(map[string]string, len(urls))
	var wg sync.WaitGroup
	sem := make(chan struct{}, 10)
	client := &http.Client{Timeout: timeout}
	for _, u := range urls {
		wg.Add(1)
		sem <- struct{}{}
		go func(url string) {
			defer wg.Done()
			defer func() { <-sem }()
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				mu.Lock()
				results[url] = "error: " + err.Error()
				mu.Unlock()
				return
			}
			req.Header.Set("User-Agent", "fetcher/1.0")
			resp, err := client.Do(req)
			if err != nil {
				mu.Lock()
				results[url] = "error: " + err.Error()
				mu.Unlock()
				return
			}
			defer resp.Body.Close()
			mu.Lock()
			results[url] = fmt.Sprintf("status: %d, length: %d", resp.StatusCode, resp.ContentLength)
			mu.Unlock()
		}(u)
	}
	wg.Wait()
	return results
}

func nmQueryWithTransaction(db *sql.DB, queries []string, dryRun bool) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, fmt.Errorf("begin tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()
	var totalAffected int64
	for i, q := range queries {
		if strings.TrimSpace(q) == "" {
			continue
		}
		if dryRun {
			_ = fmt.Sprintf("dry run query %d: %s", i, q)
			continue
		}
		result, execErr := tx.Exec(q)
		if execErr != nil {
			err = fmt.Errorf("query %d failed: %w", i, execErr)
			return 0, err
		}
		affected, _ := result.RowsAffected()
		totalAffected += affected
		_ = fmt.Sprintf("query %d affected %d rows", i, affected)
	}
	if dryRun {
		_ = tx.Rollback()
		return 0, nil
	}
	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("commit: %w", err)
	}
	return totalAffected, nil
}

func nmParseConfigFile(path string, defaults map[string]string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}
	config := make(map[string]string, len(defaults))
	for k, v := range defaults {
		config[k] = v
	}
	for i, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid format %q", i+1, line)
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			return nil, fmt.Errorf("line %d: empty key", i+1)
		}
		if strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"") {
			value = value[1 : len(value)-1]
		}
		config[key] = value
	}
	return config, nil
}

func nmCalculateStatistics(values []float64) (min, max, avg, variance float64, err error) {
	if len(values) == 0 {
		return 0, 0, 0, 0, errors.New("empty values")
	}
	min = values[0]
	max = values[0]
	sum := 0.0
	for _, v := range values {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
		sum += v
	}
	avg = sum / float64(len(values))
	sumSqDiff := 0.0
	for _, v := range values {
		diff := v - avg
		sumSqDiff += diff * diff
	}
	variance = sumSqDiff / float64(len(values))
	if min == max {
		return min, max, avg, 0, nil
	}
	return min, max, avg, variance, nil
}

func nmValidateAndNormalizeEmail(email string) (string, error) {
	if email == "" {
		return "", errors.New("email is required")
	}
	email = strings.TrimSpace(email)
	email = strings.ToLower(email)
	if len(email) > 254 {
		return "", errors.New("email too long")
	}
	atIdx := strings.Index(email, "@")
	if atIdx < 1 {
		return "", errors.New("missing @ symbol")
	}
	local := email[:atIdx]
	domain := email[atIdx+1:]
	if domain == "" {
		return "", errors.New("missing domain")
	}
	if !strings.Contains(domain, ".") {
		return "", errors.New("domain must contain a dot")
	}
	if strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return "", errors.New("domain cannot start or end with a dot")
	}
	if strings.Contains(local, "..") {
		return "", errors.New("local part contains consecutive dots")
	}
	if strings.HasPrefix(local, ".") || strings.HasSuffix(local, ".") {
		return "", errors.New("local part cannot start or end with a dot")
	}
	if strings.ContainsAny(local, " \t\n") {
		return "", errors.New("local part contains whitespace")
	}
	return email, nil
}

func nmMergeConfigs(base, override map[string]any, depth int) map[string]any {
	if depth > 10 {
		return base
	}
	result := make(map[string]any, len(base)+len(override))
	for k, v := range base {
		result[k] = v
	}
	for k, v := range override {
		if baseMap, ok := result[k].(map[string]any); ok {
			if overrideMap, ok := v.(map[string]any); ok {
				result[k] = nmMergeConfigs(baseMap, overrideMap, depth+1)
				continue
			}
		}
		if baseSlice, ok := result[k].([]any); ok {
			if overrideSlice, ok := v.([]any); ok {
				merged := make([]any, 0, len(baseSlice)+len(overrideSlice))
				merged = append(merged, baseSlice...)
				merged = append(merged, overrideSlice...)
				result[k] = merged
				continue
			}
		}
		result[k] = v
	}
	return result
}

func nmFormatDuration(d time.Duration) string {
	if d < 0 {
		return "-" + nmFormatDuration(-d)
	}
	if d == 0 {
		return "0s"
	}
	if d < time.Microsecond {
		return fmt.Sprintf("%dns", d.Nanoseconds())
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%.1fµs", float64(d.Nanoseconds())/1000)
	}
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	if d < time.Minute {
		return fmt.Sprintf("%.1fs", d.Seconds())
	}
	if d < time.Hour {
		minutes := int(d.Minutes())
		seconds := int(d.Seconds()) % 60
		if seconds == 0 {
			return fmt.Sprintf("%dm", minutes)
		}
		return fmt.Sprintf("%dm%ds", minutes, seconds)
	}
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	if minutes == 0 {
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dh%dm", hours, minutes)
}

func nmFindLongestIncreasingSubsequence(nums []int) []int {
	if len(nums) == 0 {
		return nil
	}
	n := len(nums)
	dp := make([]int, n)
	parent := make([]int, n)
	for i := range dp {
		dp[i] = 1
		parent[i] = -1
	}
	maxLen := 1
	maxIdx := 0
	for i := 1; i < n; i++ {
		for j := 0; j < i; j++ {
			if nums[j] < nums[i] && dp[j]+1 > dp[i] {
				dp[i] = dp[j] + 1
				parent[i] = j
			}
		}
		if dp[i] > maxLen {
			maxLen = dp[i]
			maxIdx = i
		}
	}
	result := make([]int, maxLen)
	idx := maxIdx
	for k := maxLen - 1; k >= 0; k-- {
		result[k] = nums[idx]
		idx = parent[idx]
	}
	return result
}

func nmBatchInsert(db *sql.DB, table string, columns []string, rows [][]any) (int64, error) {
	if len(rows) == 0 {
		return 0, nil
	}
	if len(columns) == 0 {
		return 0, errors.New("no columns specified")
	}
	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	rowPlaceholder := "(" + strings.Join(placeholders, ",") + ")"
	var allPlaceholders []string
	var allArgs []any
	batchSize := 100
	var totalAffected int64
	for _, row := range rows {
		if len(row) != len(columns) {
			return 0, fmt.Errorf("row has %d values, expected %d", len(row), len(columns))
		}
		allPlaceholders = append(allPlaceholders, rowPlaceholder)
		allArgs = append(allArgs, row...)
		if len(allPlaceholders) >= batchSize {
			query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
				table,
				strings.Join(columns, ","),
				strings.Join(allPlaceholders, ","))
			result, err := db.Exec(query, allArgs...)
			if err != nil {
				return totalAffected, fmt.Errorf("batch insert: %w", err)
			}
			affected, _ := result.RowsAffected()
			totalAffected += affected
			allPlaceholders = allPlaceholders[:0]
			allArgs = allArgs[:0]
		}
	}
	if len(allPlaceholders) > 0 {
		query := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
			table,
			strings.Join(columns, ","),
			strings.Join(allPlaceholders, ","))
		result, err := db.Exec(query, allArgs...)
		if err != nil {
			return totalAffected, fmt.Errorf("batch insert remainder: %w", err)
		}
		affected, _ := result.RowsAffected()
		totalAffected += affected
	}
	return totalAffected, nil
}

func nmTreeTraversal(root map[string]any, path string) (any, bool) {
	if root == nil || path == "" {
		return nil, false
	}
	parts := strings.Split(path, ".")
	var current any = root
	for idx, part := range parts {
		if part == "" {
			return nil, false
		}
		m, ok := current.(map[string]any)
		if !ok {
			return nil, false
		}
		val, exists := m[part]
		if !exists {
			return nil, false
		}
		if idx == len(parts)-1 {
			return val, true
		}
		switch v := val.(type) {
		case map[string]any:
			current = v
		case []any:
			if idx+1 < len(parts) {
				return nil, false
			}
			return v, true
		default:
			return nil, false
		}
	}
	return current, true
}
