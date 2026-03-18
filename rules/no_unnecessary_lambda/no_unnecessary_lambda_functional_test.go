package no_unnecessary_lambda_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_unnecessary_lambda"
)

func TestNoUnnecessaryLambda(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_unnecessary_lambda", &no_unnecessary_lambda.NoUnnecessaryLambdaRule{})
}
