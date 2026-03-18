package no_todo_without_detail_test

import (
	"testing"

	"github.com/strowk/vint/rules/functional_test_helpers"
	"github.com/strowk/vint/rules/no_todo_without_detail"
)

func TestNoTodoWithoutDetail(t *testing.T) {
	functional_test_helpers.TestRule(t, "no_todo_without_detail", &no_todo_without_detail.NoTodoWithoutDetailRule{})
}
