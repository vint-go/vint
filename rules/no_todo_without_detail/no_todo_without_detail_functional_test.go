package no_todo_without_detail_test

import (
	"testing"

	"github.com/vint-go/vint/rules/functional_test_helpers"
	"github.com/vint-go/vint/rules/no_todo_without_detail"
)

func TestNoTodoWithoutDetail(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_todo_without_detail", &no_todo_without_detail.NoTodoWithoutDetailRule{})
}
