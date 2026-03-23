package no_http_client_post_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_http_client_post"
)

func TestNoHttpClientPost(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_http_client_post", &no_http_client_post.NoHttpClientPostRule{})
}
