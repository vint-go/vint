package migrate

// This file serves as a registry of configuration settings that are
// not supported during migration from golangci-lint to vint.
// Each function returns a human-readable warning string.

// --- gosec ---

// WarnGosecSeverityFilter warns that gosec's severity filter has no equivalent in vint.
//
// In gosec, severity filters findings by their severity level (low/medium/high),
// suppressing lower-severity results before they are reported. Vint has severity
// in its data model (rules can be configured with a severity level), but it does
// not support filtering findings by severity at the level of a single linter's
// settings. Users who relied on severity filtering to suppress low-priority gosec
// findings should disable the specific rules they consider unnecessary.
//
// TODO: This could be supported by mapping each gosec rule ID to its known
// severity level and excluding rules below the user's threshold during migration,
// effectively translating the severity filter into a set of disabled rules.
func WarnGosecSeverityFilter() string {
	return `gosec: "severity" filter is not supported in vint — all matching findings are reported`
}

// WarnGosecConfidenceFilter warns that gosec's confidence filter has no equivalent in vint.
//
// In gosec, the confidence filter suppresses findings below a confidence threshold
// (low/medium/high), reducing false positives at the cost of missed issues. Vint
// rules are designed to operate at a single confidence level with low false-positive
// rates, so a confidence dial is unnecessary. Users seeing excessive noise from a
// specific rule should disable it.
//
// TODO: Same as severity — each gosec rule has a known confidence level, so the
// migrator could exclude rules below the user's confidence threshold.
func WarnGosecConfidenceFilter() string {
	return `gosec: "confidence" filter is not supported in vint — all matching findings are reported`
}

// WarnGosecGlobalNosec warns that gosec's global.nosec setting is not migrated.
//
// In gosec, global.nosec controls whether #nosec annotations in source code are
// honored. When set to true, all #nosec annotations are ignored and every finding
// is reported. Vint uses //nolint: directives for suppression, which are always
// honored. There is no global switch to disable them because suppression control
// is handled per-directive by the nolintlint rule instead.
func WarnGosecGlobalNosec() string {
	return `gosec: "global.nosec" setting is not supported in vint — use "//nolint:" directives instead`
}

// WarnGosecGlobalAudit warns that gosec's global.audit setting is not migrated.
//
// In gosec, audit mode enables additional stricter checks and lowers thresholds
// across all rules. Vint rules each run at a single strictness level and do not
// have a global "be stricter" toggle. Users who need tighter checks should
// configure individual rule thresholds (e.g. lower entropy thresholds, stricter
// permission modes) rather than relying on a blanket audit flag.
func WarnGosecGlobalAudit() string {
	return `gosec: "global.audit" setting is not supported in vint — rules run at their default strictness`
}

// WarnGosecG104Config warns that gosec G104 per-rule config (unchecked errors) is not migrated.
//
// In gosec, G104's config specifies a list of packages/functions whose unchecked
// error returns should be flagged. This is an include-list ("only check these").
// Vint's noUncheckedError rule uses the opposite model — an exclude-list of
// functions whose errors are safe to ignore. The semantics are inverted and cannot
// be mechanically translated: gosec's include-list means "ignore everything except
// these", while vint's exclude-list means "check everything except these".
func WarnGosecG104Config() string {
	return "gosec: per-rule config for G104 (unchecked errors) is not supported in vint"
}

// WarnGosecG111Config warns that gosec G111 per-rule config (directory serving pattern) is not migrated.
//
// In gosec, G111's pattern config allows a custom regex to detect which paths are
// considered risky when passed to http.Dir(). Vint's noFilesystemRootServing rule
// uses a simpler approach: it only flags http.Dir("/") — serving the literal
// filesystem root. Supporting arbitrary path patterns would require a substantially
// different detection model with higher false-positive risk, so we check only the
// most dangerous case.
func WarnGosecG111Config() string {
	return "gosec: per-rule config for G111 (directory serving pattern) is not supported in vint"
}
