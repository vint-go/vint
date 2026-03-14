package rulecache

// CachedFailure is the serializable form of a lint failure.
// It mirrors lint.Failure but flattens positions and drops the ast.Node.
type CachedFailure struct {
	Message         string
	RuleName        string
	Category        string
	StartFilename   string
	StartLine       int
	StartColumn     int
	EndFilename     string
	EndLine         int
	EndColumn       int
	Confidence      float64
	ReplacementLine string
}
