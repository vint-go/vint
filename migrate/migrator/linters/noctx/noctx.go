package noctx

import "github.com/vint-go/vint/migrate"

func init() {
	migrate.RegisterMigrator(&Migrator{})
}

// Migrator handles migration for the noctx golangci-lint linter.
// noctx has no configuration options — it checks that functions
// accepting a context.Context are called with the context-aware variant.
type Migrator struct{}

func (*Migrator) Name() string {
	return "noctx"
}

func (*Migrator) MigrateConfig(_ map[string]any) (map[string]migrate.VintRuleConfig, error) {
	return map[string]migrate.VintRuleConfig{
		"lint/correctness/noExecCommand":       {},
		"lint/correctness/noHttpClientGet":      {},
		"lint/correctness/noHttpClientHead":     {},
		"lint/correctness/noHttpClientPost":     {},
		"lint/correctness/noHttpClientPostForm": {},
		"lint/correctness/noHttpGet":            {},
		"lint/correctness/noHttpHead":           {},
		"lint/correctness/noHttpNewRequest":     {},
		"lint/correctness/noHttpPost":           {},
		"lint/correctness/noHttpPostForm":       {},
		"lint/correctness/noHttptestNewRequest": {},
		"lint/correctness/noNetDial":            {},
		"lint/correctness/noNetDialTimeout":     {},
		"lint/correctness/noNetListen":          {},
		"lint/correctness/noNetListenPacket":    {},
		"lint/correctness/noNetLookupAddr":      {},
		"lint/correctness/noNetLookupCname":     {},
		"lint/correctness/noNetLookupHost":      {},
		"lint/correctness/noNetLookupIp":        {},
		"lint/correctness/noNetLookupMx":        {},
		"lint/correctness/noNetLookupNs":        {},
		"lint/correctness/noNetLookupPort":      {},
		"lint/correctness/noNetLookupSrv":       {},
		"lint/correctness/noNetLookupTxt":       {},
		"lint/correctness/noSqlDbBegin":         {},
		"lint/correctness/noSqlDbExec":          {},
		"lint/correctness/noSqlDbPing":          {},
		"lint/correctness/noSqlDbPrepare":       {},
		"lint/correctness/noSqlDbQuery":         {},
		"lint/correctness/noSqlDbQueryRow":      {},
		"lint/correctness/noSqlStmtExec":        {},
		"lint/correctness/noSqlStmtQuery":       {},
		"lint/correctness/noSqlStmtQueryRow":    {},
		"lint/correctness/noSqlTxExec":          {},
		"lint/correctness/noSqlTxPrepare":       {},
		"lint/correctness/noSqlTxQuery":         {},
		"lint/correctness/noSqlTxQueryRow":      {},
		"lint/correctness/noSqlTxStmt":          {},
		"lint/correctness/noTlsConnHandshake":   {},
		"lint/correctness/noTlsDial":            {},
		"lint/correctness/noTlsDialWithDialer":  {},
	}, nil
}
