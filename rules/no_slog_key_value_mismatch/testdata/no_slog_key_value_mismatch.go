package fixtures

import "log/slog"

// Invalid: odd number of arguments, "status" has no value
func badOddArgs() {
	slog.Info("request", "method", "GET", "status") // MATCH /slog key-value mismatch: missing value for key argument/
}

// Invalid: non-string key
func badNonStringKey() {
	slog.Info("request", 200, "OK") // MATCH /slog key-value mismatch: non-string key argument/
}

// Valid: properly paired key-value arguments
func goodPaired() {
	slog.Info("request", "method", "GET", "status", 200)
}

// Valid: using slog.Attr for typed attributes
func goodAttr() {
	slog.Info("request",
		slog.String("method", "GET"),
		slog.Int("status", 200),
	)
}

// Valid: no extra args (just message)
func goodNoArgs() {
	slog.Info("request")
}

// Valid: mixed slog.Attr and key-value pairs
func goodMixed() {
	slog.Info("request",
		slog.String("method", "GET"),
		"status", 200,
	)
}

// Invalid: odd args with Error
func badErrorOddArgs() {
	slog.Error("failed", "code") // MATCH /slog key-value mismatch: missing value for key argument/
}

// Valid: Debug with correct pairs
func goodDebug() {
	slog.Debug("debug", "key", "value")
}

// Valid: Warn with correct pairs
func goodWarn() {
	slog.Warn("warning", "key", "value")
}
