package lint

import (
	"go/token"

	"github.com/strowk/vint/internal/rulecache"
)

// DisabledInterval contains a single disabled interval and the associated rule name.
type DisabledInterval struct {
	From     token.Position
	To       token.Position
	RuleName string
}

// Rule defines an abstract rule interface.
type Rule interface {
	Name() string
	Apply(*File, Arguments) []Failure
}

// Grouped defines an interface for rules that belong to a group.
type Grouped interface {
	Group() string
}

// ConfigurableRule defines an abstract configurable rule interface.
type ConfigurableRule interface {
	Configure(Arguments) error
}

// FullRuleName returns the fully qualified rule name.
// For rules implementing Grouped, it returns "lint/<group>/<name>".
// For other rules, it returns the rule's Name().
func FullRuleName(r Rule) string {
	if g, ok := r.(Grouped); ok {
		return "lint/" + g.Group() + "/" + r.Name()
	}
	return r.Name()
}

// TieredRule is implemented by rules that declare their cache tier.
// Rules that do not implement this default to TierCrossPackage (safest).
type TieredRule interface {
	Rule
	CacheTier() rulecache.CacheTier
}

// UncacheableRule is implemented by rules whose results must never be cached.
type UncacheableRule interface {
	Uncacheable() bool
}

// ToFailurePosition returns the failure position.
func ToFailurePosition(start, end token.Pos, file *File) FailurePosition {
	return FailurePosition{
		Start: file.ToPosition(start),
		End:   file.ToPosition(end),
	}
}
