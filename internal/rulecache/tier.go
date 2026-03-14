package rulecache

// CacheTier declares how much context a rule needs beyond the file itself.
type CacheTier int

const (
	// TierFileOnly means the rule only inspects the file's AST and content.
	TierFileOnly CacheTier = iota
	// TierPackageAware means the rule also uses type info from the same package.
	TierPackageAware
	// TierCrossPackage means the rule uses type info from imported packages.
	TierCrossPackage
)
