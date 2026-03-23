package no_http_post_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_http_post"
)

func TestNoHttpPost(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_http_post", &no_http_post.NoHttpPostRule{})
}
