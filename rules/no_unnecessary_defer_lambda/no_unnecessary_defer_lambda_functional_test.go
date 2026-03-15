package no_unnecessary_defer_lambda_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unnecessary_defer_lambda"
)

func TestNoUnnecessaryDeferLambda(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_defer_lambda", &no_unnecessary_defer_lambda.NoUnnecessaryDeferLambdaRule{})
}
