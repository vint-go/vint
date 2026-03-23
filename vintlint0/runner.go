package vintlint0

import (
	"errors"

	"github.com/vint-go/vint/lint"
)

// runRulesOnFile executes the given rules against a single file, sending
// results to the failures channel. Respects disabled intervals and confidence
// thresholds. Caching is handled at the package level by the caller.
func runRulesOnFile(file *lint.File, rules []lint.Rule, config lint.Config, failures chan<- lint.Failure) error {
	rulesConfig := config.Rules

	// Compute disabled intervals from vint-ignore directives in comments.
	disabledIntervals := parseVintIgnoreSuppressions(file, failures)

	for _, currentRule := range rules {
		fullName := lint.FullRuleName(currentRule)
		ruleConfig := rulesConfig[fullName]
		if ruleConfig.MustExclude(file.Name) {
			continue
		}

		// Run the rule.
		currentFailures := currentRule.Apply(file, ruleConfig.Arguments)
		for idx, failure := range currentFailures {
			if failure.IsInternal() {
				return errors.New(failure.Failure)
			}

			if failure.RuleName == "" {
				failure.RuleName = fullName
			}
			if failure.Node != nil {
				failure.Position = lint.ToFailurePosition(failure.Node.Pos(), failure.Node.End(), file)
			}
			currentFailures[idx] = failure
		}

		// Filter by disabled intervals.
		currentFailures = lint.FilterFailuresByDisabledIntervals(currentFailures, disabledIntervals)

		// Send passing failures to the channel.
		for _, failure := range currentFailures {
			if failure.Confidence >= config.Confidence {
				failures <- failure
			}
		}
	}

	return nil
}
