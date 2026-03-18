package fixtures

import (
	"log/slog"
	"strings"
)

// Invalid: odd number of arguments to strings.NewReplacer
func badNewReplacer() {
	strings.NewReplacer("old", "new", "extra") // MATCH /odd number of arguments passed to NewReplacer, expected even number (key-value pairs)/
}

// Invalid: odd number of arguments to slog.Info (missing value for "key2")
func badSlogInfo() {
	slog.Info("message", "key1", "value1", "key2") // MATCH /odd number of arguments passed to slog function, expected even number (key-value pairs)/
}

// Invalid: odd number of arguments to slog.Error
func badSlogError() {
	slog.Error("failed", "code") // MATCH /odd number of arguments passed to slog function, expected even number (key-value pairs)/
}

// Invalid: odd number of arguments to slog.Debug
func badSlogDebug() {
	slog.Debug("debug", "key") // MATCH /odd number of arguments passed to slog function, expected even number (key-value pairs)/
}

// Invalid: odd number of arguments to slog.Warn
func badSlogWarn() {
	slog.Warn("warning", "key") // MATCH /odd number of arguments passed to slog function, expected even number (key-value pairs)/
}

// Valid: even number of arguments to strings.NewReplacer
func goodNewReplacer() {
	strings.NewReplacer("old", "new", "old2", "new2")
}

// Valid: single pair to strings.NewReplacer
func goodNewReplacerSingle() {
	strings.NewReplacer("old", "new")
}

// Valid: even number of arguments to slog.Info
func goodSlogInfo() {
	slog.Info("message", "key1", "value1", "key2", "value2")
}

// Valid: no extra args (just message)
func goodSlogInfoNoArgs() {
	slog.Info("message")
}

// Valid: using slog.Attr for typed attributes
func goodSlogInfoAttr() {
	slog.Info("message",
		slog.String("key1", "value1"),
		slog.Int("key2", 42),
	)
}

// Valid: mixed slog.Attr and key-value pairs
func goodSlogInfoMixed() {
	slog.Info("message",
		slog.String("key1", "value1"),
		"key2", "value2",
	)
}

// Valid: empty call to strings.NewReplacer
func goodNewReplacerEmpty() {
	strings.NewReplacer()
}
