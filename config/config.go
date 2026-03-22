// Package config implements revive's configuration data structures and related methods.
package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"

	"github.com/strowk/vint/formatter"
	"github.com/strowk/vint/lint"
	"github.com/strowk/vint/rule"
	"github.com/strowk/vint/rules/no_address_nil_comparison"
	"github.com/strowk/vint/rules/no_address_of_dereference"
	"github.com/strowk/vint/rules/no_always_true_len_check"
	"github.com/strowk/vint/rules/no_asm_decl_mismatch"
	"github.com/strowk/vint/rules/no_atoi_overflow"
	"github.com/strowk/vint/rules/no_atomic_alignment_issue"
	"github.com/strowk/vint/rules/no_atomic_assign_misuse"
	"github.com/strowk/vint/rules/no_bad_lock_pattern"
	"github.com/strowk/vint/rules/no_bad_regexp_pattern"
	"github.com/strowk/vint/rules/no_bad_sort_usage"
	"github.com/strowk/vint/rules/no_banned_characters"
	"github.com/strowk/vint/rules/no_benchmark_n_assignment"
	"github.com/strowk/vint/rules/no_bind_to_all_interfaces"
	"github.com/strowk/vint/rules/no_blank_error_assignment"
	"github.com/strowk/vint/rules/no_blank_import"
	"github.com/strowk/vint/rules/no_builtin_shadow"
	"github.com/strowk/vint/rules/no_builtin_shadow_decl"
	"github.com/strowk/vint/rules/no_call_to_gc"
	"github.com/strowk/vint/rules/no_capitalized_error_string"
	"github.com/strowk/vint/rules/no_capitalized_local"
	"github.com/strowk/vint/rules/no_cgi_import"
	"github.com/strowk/vint/rules/no_cgo_pointer_violation"
	"github.com/strowk/vint/rules/no_command_injection_taint"
	"github.com/strowk/vint/rules/no_commented_out_code"
	"github.com/strowk/vint/rules/no_commented_out_import"
	"github.com/strowk/vint/rules/no_conflicting_http_mux_patterns"
	"github.com/strowk/vint/rules/no_confusing_naming"
	"github.com/strowk/vint/rules/no_confusing_results"
	"github.com/strowk/vint/rules/no_constant_parameter"
	"github.com/strowk/vint/rules/no_constant_result"
	"github.com/strowk/vint/rules/no_context_keys_type"
	"github.com/strowk/vint/rules/no_context_propagation_failure"
	"github.com/strowk/vint/rules/no_control_char_in_string"
	"github.com/strowk/vint/rules/no_copied_lock"
	"github.com/strowk/vint/rules/no_deep_equal_errors"
	"github.com/strowk/vint/rules/no_deep_exit"
	"github.com/strowk/vint/rules/no_default_slice_index"
	"github.com/strowk/vint/rules/no_defer_close_before_err_check"
	"github.com/strowk/vint/rules/no_defer_gotcha"
	"github.com/strowk/vint/rules/no_defer_in_infinite_loop"
	"github.com/strowk/vint/rules/no_defer_in_loop"
	"github.com/strowk/vint/rules/no_defer_time_misuse"
	"github.com/strowk/vint/rules/no_denied_import"
	"github.com/strowk/vint/rules/no_deprecated_hash_function"
	"github.com/strowk/vint/rules/no_deprecated_usage"
	"github.com/strowk/vint/rules/no_direct_error_comparison"
	"github.com/strowk/vint/rules/no_discarded_append"
	"github.com/strowk/vint/rules/no_doc_comment_stub"
	"github.com/strowk/vint/rules/no_dot_import"
	"github.com/strowk/vint/rules/no_dubious_bit_shift"
	"github.com/strowk/vint/rules/no_duplicate_argument"
	"github.com/strowk/vint/rules/no_duplicate_branch_body"
	"github.com/strowk/vint/rules/no_duplicate_build_constraint"
	"github.com/strowk/vint/rules/no_duplicate_case"
	"github.com/strowk/vint/rules/no_duplicate_code"
	"github.com/strowk/vint/rules/no_duplicate_constants"
	"github.com/strowk/vint/rules/no_duplicate_cutset_chars"
	"github.com/strowk/vint/rules/no_duplicate_if_condition"
	"github.com/strowk/vint/rules/no_duplicate_import"
	"github.com/strowk/vint/rules/no_duplicate_option"
	"github.com/strowk/vint/rules/no_duplicate_sub_expression"
	"github.com/strowk/vint/rules/no_duplicated_imports"
	"github.com/strowk/vint/rules/no_dynamic_errors"
	"github.com/strowk/vint/rules/no_dynamic_format_string"
	"github.com/strowk/vint/rules/no_empty_branch"
	"github.com/strowk/vint/rules/no_empty_critical_section"
	"github.com/strowk/vint/rules/no_empty_declaration"
	"github.com/strowk/vint/rules/no_empty_fallthrough"
	"github.com/strowk/vint/rules/no_empty_for_loop"
	"github.com/strowk/vint/rules/no_error_strings"
	"github.com/strowk/vint/rules/no_eval_order_dependency"
	"github.com/strowk/vint/rules/no_excessive_arguments"
	"github.com/strowk/vint/rules/no_excessive_blank_identifiers"
	"github.com/strowk/vint/rules/no_excessive_control_nesting"
	"github.com/strowk/vint/rules/no_excessive_file_length"
	"github.com/strowk/vint/rules/no_excessive_function_results"
	"github.com/strowk/vint/rules/no_excessive_public_structs"
	"github.com/strowk/vint/rules/no_excessive_results"
	"github.com/strowk/vint/rules/no_excessive_shift"
	"github.com/strowk/vint/rules/no_excessive_statements"
	"github.com/strowk/vint/rules/no_exec_command"
	"github.com/strowk/vint/rules/no_exit_after_defer"
	"github.com/strowk/vint/rules/no_explicit_bool_comparison"
	"github.com/strowk/vint/rules/no_exposed_pprof"
	"github.com/strowk/vint/rules/no_exposed_sync_mutex"
	"github.com/strowk/vint/rules/no_external_error_reassign"
	"github.com/strowk/vint/rules/no_file_inclusion_via_variable"
	"github.com/strowk/vint/rules/no_file_scoped_denied_import"
	"github.com/strowk/vint/rules/no_filesystem_root_serving"
	"github.com/strowk/vint/rules/no_filesystem_toctou"
	"github.com/strowk/vint/rules/no_flag_deref_before_parse"
	"github.com/strowk/vint/rules/no_flag_parameter"
	"github.com/strowk/vint/rules/no_forbidden_call_in_wg_go"
	"github.com/strowk/vint/rules/no_frame_pointer_clobber"
	"github.com/strowk/vint/rules/no_guard_around_delete"
	"github.com/strowk/vint/rules/no_guard_around_map_access"
	"github.com/strowk/vint/rules/no_hardcoded_credentials"
	"github.com/strowk/vint/rules/no_hardcoded_iv"
	"github.com/strowk/vint/rules/no_high_cognitive_complexity"
	"github.com/strowk/vint/rules/no_high_cyclomatic_complexity"
	"github.com/strowk/vint/rules/no_http_client_get"
	"github.com/strowk/vint/rules/no_http_client_head"
	"github.com/strowk/vint/rules/no_http_client_post"
	"github.com/strowk/vint/rules/no_http_client_post_form"
	"github.com/strowk/vint/rules/no_http_get"
	"github.com/strowk/vint/rules/no_http_head"
	"github.com/strowk/vint/rules/no_http_new_request"
	"github.com/strowk/vint/rules/no_http_post"
	"github.com/strowk/vint/rules/no_http_post_form"
	"github.com/strowk/vint/rules/no_http_request_smuggling"
	"github.com/strowk/vint/rules/no_http_response_misuse"
	"github.com/strowk/vint/rules/no_httptest_new_request"
	"github.com/strowk/vint/rules/no_huge_param"
	"github.com/strowk/vint/rules/no_identical_branches"
	"github.com/strowk/vint/rules/no_identical_if_else_if_branches"
	"github.com/strowk/vint/rules/no_identical_if_else_if_conditions"
	"github.com/strowk/vint/rules/no_identical_switch_branches"
	"github.com/strowk/vint/rules/no_identical_switch_conditions"
	"github.com/strowk/vint/rules/no_ignored_query_modification"
	"github.com/strowk/vint/rules/no_ignored_query_result"
	"github.com/strowk/vint/rules/no_immediate_new_deref"
	"github.com/strowk/vint/rules/no_implicit_const_value"
	"github.com/strowk/vint/rules/no_import_shadow"
	"github.com/strowk/vint/rules/no_impossible_builtin_result"
	"github.com/strowk/vint/rules/no_impossible_condition"
	"github.com/strowk/vint/rules/no_impossible_interface_assert"
	"github.com/strowk/vint/rules/no_impossible_nil_comparison"
	"github.com/strowk/vint/rules/no_imprecise_constant"
	"github.com/strowk/vint/rules/no_inappropriate_context_key"
	"github.com/strowk/vint/rules/no_incorrect_time_format"
	"github.com/strowk/vint/rules/no_ineffective_bitwise_op"
	"github.com/strowk/vint/rules/no_ineffective_break"
	"github.com/strowk/vint/rules/no_ineffective_rand_call"
	"github.com/strowk/vint/rules/no_ineffectual_assignment"
	"github.com/strowk/vint/rules/no_inefficient_map_lookup"
	"github.com/strowk/vint/rules/no_infinite_recursion"
	"github.com/strowk/vint/rules/no_init_function"
	"github.com/strowk/vint/rules/no_inline_sync_once_func"
	"github.com/strowk/vint/rules/no_insecure_cookie"
	"github.com/strowk/vint/rules/no_insecure_host_key_callback"
	"github.com/strowk/vint/rules/no_insecure_random"
	"github.com/strowk/vint/rules/no_insecure_tls_config"
	"github.com/strowk/vint/rules/no_integer_division_truncation"
	"github.com/strowk/vint/rules/no_integer_overflow_conversion"
	"github.com/strowk/vint/rules/no_invalid_binary_arg"
	"github.com/strowk/vint/rules/no_invalid_errors_as"
	"github.com/strowk/vint/rules/no_invalid_exec_command_arg"
	"github.com/strowk/vint/rules/no_invalid_flag_name"
	"github.com/strowk/vint/rules/no_invalid_host_port"
	"github.com/strowk/vint/rules/no_invalid_regexp"
	"github.com/strowk/vint/rules/no_invalid_sort_slice_arg"
	"github.com/strowk/vint/rules/no_invalid_strconv_arg"
	"github.com/strowk/vint/rules/no_invalid_template"
	"github.com/strowk/vint/rules/no_invalid_unsafe_pointer"
	"github.com/strowk/vint/rules/no_invalid_url_parse"
	"github.com/strowk/vint/rules/no_invalid_utf8_string_arg"
	"github.com/strowk/vint/rules/no_invariant_loop_condition"
	"github.com/strowk/vint/rules/no_leading_blank_line"
	"github.com/strowk/vint/rules/no_line_too_long"
	"github.com/strowk/vint/rules/no_log_injection_taint"
	"github.com/strowk/vint/rules/no_long_functions"
	"github.com/strowk/vint/rules/no_loop_closure_capture"
	"github.com/strowk/vint/rules/no_lost_cancel"
	"github.com/strowk/vint/rules/no_magic_number_in_argument"
	"github.com/strowk/vint/rules/no_magic_number_in_assignment"
	"github.com/strowk/vint/rules/no_magic_number_in_case"
	"github.com/strowk/vint/rules/no_magic_number_in_condition"
	"github.com/strowk/vint/rules/no_magic_number_in_operation"
	"github.com/strowk/vint/rules/no_magic_number_in_return"
	"github.com/strowk/vint/rules/no_malformed_build_tag"
	"github.com/strowk/vint/rules/no_malformed_directive"
	"github.com/strowk/vint/rules/no_malformed_struct_tag"
	"github.com/strowk/vint/rules/no_malformed_test_function"
	"github.com/strowk/vint/rules/no_mismatched_append_assign"
	"github.com/strowk/vint/rules/no_misplaced_default_case"
	"github.com/strowk/vint/rules/no_missing_read_header_timeout"
	"github.com/strowk/vint/rules/no_missing_return_after_http_error"
	"github.com/strowk/vint/rules/no_misspelled_words"
	"github.com/strowk/vint/rules/no_mixed_case_hex_literal"
	"github.com/strowk/vint/rules/no_modified_parameter"
	"github.com/strowk/vint/rules/no_modified_value_receiver"
	"github.com/strowk/vint/rules/no_modulo_one"
	"github.com/strowk/vint/rules/no_multi_line_func_break"
	"github.com/strowk/vint/rules/no_multi_line_if_break"
	"github.com/strowk/vint/rules/no_naked_return"
	"github.com/strowk/vint/rules/no_nan_comparison"
	"github.com/strowk/vint/rules/no_nested_structs"
	"github.com/strowk/vint/rules/no_net_dial"
	"github.com/strowk/vint/rules/no_net_dial_timeout"
	"github.com/strowk/vint/rules/no_net_ip_bytes_equal"
	"github.com/strowk/vint/rules/no_net_listen"
	"github.com/strowk/vint/rules/no_net_listen_packet"
	"github.com/strowk/vint/rules/no_net_lookup_addr"
	"github.com/strowk/vint/rules/no_net_lookup_cname"
	"github.com/strowk/vint/rules/no_net_lookup_host"
	"github.com/strowk/vint/rules/no_net_lookup_ip"
	"github.com/strowk/vint/rules/no_net_lookup_mx"
	"github.com/strowk/vint/rules/no_net_lookup_ns"
	"github.com/strowk/vint/rules/no_net_lookup_port"
	"github.com/strowk/vint/rules/no_net_lookup_srv"
	"github.com/strowk/vint/rules/no_net_lookup_txt"
	"github.com/strowk/vint/rules/no_never_nil_check"
	"github.com/strowk/vint/rules/no_nil_context"
	"github.com/strowk/vint/rules/no_nil_dereference"
	"github.com/strowk/vint/rules/no_nil_func_comparison"
	"github.com/strowk/vint/rules/no_nil_map_assignment"
	"github.com/strowk/vint/rules/no_nil_variable_return"
	"github.com/strowk/vint/rules/no_nolint_non_machine_readable"
	"github.com/strowk/vint/rules/no_nolint_parse_error"
	"github.com/strowk/vint/rules/no_nolint_without_explanation"
	"github.com/strowk/vint/rules/no_nolint_without_specific_linter"
	"github.com/strowk/vint/rules/no_non_canonical_header_key"
	"github.com/strowk/vint/rules/no_non_octal_file_mode"
	"github.com/strowk/vint/rules/no_non_pointer_unmarshal"
	"github.com/strowk/vint/rules/no_non_wrapping_errorf"
	"github.com/strowk/vint/rules/no_noop_function_call"
	"github.com/strowk/vint/rules/no_odd_size_slice_arg"
	"github.com/strowk/vint/rules/no_off_by_one_error"
	"github.com/strowk/vint/rules/no_overlapping_encoder_slice"
	"github.com/strowk/vint/rules/no_overwritten_argument"
	"github.com/strowk/vint/rules/no_package_directory_mismatch"
	"github.com/strowk/vint/rules/no_path_traversal_taint"
	"github.com/strowk/vint/rules/no_permissive_directory_permissions"
	"github.com/strowk/vint/rules/no_permissive_file_permissions"
	"github.com/strowk/vint/rules/no_permissive_os_create"
	"github.com/strowk/vint/rules/no_permissive_write_file_permissions"
	"github.com/strowk/vint/rules/no_pointer_to_ref_param"
	"github.com/strowk/vint/rules/no_predictable_temp_file"
	"github.com/strowk/vint/rules/no_printf_format_mismatch"
	"github.com/strowk/vint/rules/no_range_append_all"
	"github.com/strowk/vint/rules/no_range_expr_copy"
	"github.com/strowk/vint/rules/no_range_val_copy"
	"github.com/strowk/vint/rules/no_range_variable_alias"
	"github.com/strowk/vint/rules/no_redundant_build_tag"
	"github.com/strowk/vint/rules/no_redundant_canonical_header_key"
	"github.com/strowk/vint/rules/no_redundant_control_flow"
	"github.com/strowk/vint/rules/no_redundant_import_alias"
	"github.com/strowk/vint/rules/no_redundant_label"
	"github.com/strowk/vint/rules/no_redundant_make_args"
	"github.com/strowk/vint/rules/no_redundant_nil_loop_check"
	"github.com/strowk/vint/rules/no_redundant_nil_slice_check"
	"github.com/strowk/vint/rules/no_redundant_nil_type_check"
	"github.com/strowk/vint/rules/no_redundant_range_val"
	"github.com/strowk/vint/rules/no_redundant_rune_conversion"
	"github.com/strowk/vint/rules/no_redundant_slice_expression"
	"github.com/strowk/vint/rules/no_redundant_sprint"
	"github.com/strowk/vint/rules/no_redundant_string_byte_conversion"
	"github.com/strowk/vint/rules/no_redundant_string_concat"
	"github.com/strowk/vint/rules/no_redundant_switch_true"
	"github.com/strowk/vint/rules/no_redundant_test_main_exit"
	"github.com/strowk/vint/rules/no_redundant_type_assertion"
	"github.com/strowk/vint/rules/no_redundant_var_type"
	"github.com/strowk/vint/rules/no_repeated_numbers"
	"github.com/strowk/vint/rules/no_repeated_strings"
	"github.com/strowk/vint/rules/no_secret_in_serialization"
	"github.com/strowk/vint/rules/no_select_break_confusion"
	"github.com/strowk/vint/rules/no_self_assignment"
	"github.com/strowk/vint/rules/no_self_referencing_finalizer"
	"github.com/strowk/vint/rules/no_separator_in_filepath_join"
	"github.com/strowk/vint/rules/no_serve_without_timeout"
	"github.com/strowk/vint/rules/no_short_rsa_key"
	"github.com/strowk/vint/rules/no_side_effect_in_init_clause"
	"github.com/strowk/vint/rules/no_single_arg_append"
	"github.com/strowk/vint/rules/no_single_case_select"
	"github.com/strowk/vint/rules/no_single_case_switch"
	"github.com/strowk/vint/rules/no_single_iteration_loop"
	"github.com/strowk/vint/rules/no_slice_bounds_out_of_range"
	"github.com/strowk/vint/rules/no_slog_key_value_mismatch"
	"github.com/strowk/vint/rules/no_smtp_injection_taint"
	"github.com/strowk/vint/rules/no_space_in_directive"
	"github.com/strowk/vint/rules/no_specific_function_call"
	"github.com/strowk/vint/rules/no_sql_concatenation"
	"github.com/strowk/vint/rules/no_sql_db_begin"
	"github.com/strowk/vint/rules/no_sql_db_exec"
	"github.com/strowk/vint/rules/no_sql_db_ping"
	"github.com/strowk/vint/rules/no_sql_db_prepare"
	"github.com/strowk/vint/rules/no_sql_db_query"
	"github.com/strowk/vint/rules/no_sql_db_query_row"
	"github.com/strowk/vint/rules/no_sql_format_string"
	"github.com/strowk/vint/rules/no_sql_injection_taint"
	"github.com/strowk/vint/rules/no_sql_stmt_exec"
	"github.com/strowk/vint/rules/no_sql_stmt_query"
	"github.com/strowk/vint/rules/no_sql_stmt_query_row"
	"github.com/strowk/vint/rules/no_sql_tx_exec"
	"github.com/strowk/vint/rules/no_sql_tx_prepare"
	"github.com/strowk/vint/rules/no_sql_tx_query"
	"github.com/strowk/vint/rules/no_sql_tx_query_row"
	"github.com/strowk/vint/rules/no_sql_tx_stmt"
	"github.com/strowk/vint/rules/no_ssh_auth_bypass"
	"github.com/strowk/vint/rules/no_ssrf_taint"
	"github.com/strowk/vint/rules/no_ssrf_via_variable"
	"github.com/strowk/vint/rules/no_std_method_signature_mismatch"
	"github.com/strowk/vint/rules/no_stdlib_version_mismatch"
	"github.com/strowk/vint/rules/no_string_index_allocation"
	"github.com/strowk/vint/rules/no_string_int_conversion"
	"github.com/strowk/vint/rules/no_strings_compare"
	"github.com/strowk/vint/rules/no_superfluous_else"
	"github.com/strowk/vint/rules/no_suspicious_map_key"
	"github.com/strowk/vint/rules/no_suspicious_sort_slice"
	"github.com/strowk/vint/rules/no_suspicious_time_sleep"
	"github.com/strowk/vint/rules/no_swapped_arguments"
	"github.com/strowk/vint/rules/no_temp_dir_deletion"
	"github.com/strowk/vint/rules/no_template_injection"
	"github.com/strowk/vint/rules/no_test_fatal_in_goroutine"
	"github.com/strowk/vint/rules/no_test_main_without_exit"
	"github.com/strowk/vint/rules/no_time_tick"
	"github.com/strowk/vint/rules/no_timer_reset_retval"
	"github.com/strowk/vint/rules/no_tls_conn_handshake"
	"github.com/strowk/vint/rules/no_tls_dial"
	"github.com/strowk/vint/rules/no_tls_dial_with_dialer"
	"github.com/strowk/vint/rules/no_tls_session_resumption_bypass"
	"github.com/strowk/vint/rules/no_todo_without_detail"
	"github.com/strowk/vint/rules/no_trailing_blank_line"
	"github.com/strowk/vint/rules/no_trojan_source_bidi"
	"github.com/strowk/vint/rules/no_truncating_comparison"
	"github.com/strowk/vint/rules/no_type_assert_else_misread"
	"github.com/strowk/vint/rules/no_unallowed_import"
	"github.com/strowk/vint/rules/no_unbounded_decompression"
	"github.com/strowk/vint/rules/no_unbounded_form_parsing"
	"github.com/strowk/vint/rules/no_unbuffered_signal_channel"
	"github.com/strowk/vint/rules/no_unchecked_error"
	"github.com/strowk/vint/rules/no_unchecked_inline_error"
	"github.com/strowk/vint/rules/no_unchecked_type_assertion"
	"github.com/strowk/vint/rules/no_unclosed_bodies"
	"github.com/strowk/vint/rules/no_unconditional_recursion"
	"github.com/strowk/vint/rules/no_unescaped_html_template"
	"github.com/strowk/vint/rules/no_unescaped_regexp_dot"
	"github.com/strowk/vint/rules/no_unexported_naming"
	"github.com/strowk/vint/rules/no_unexported_return"
	"github.com/strowk/vint/rules/no_unkeyed_literal"
	"github.com/strowk/vint/rules/no_unmarshalable_struct"
	"github.com/strowk/vint/rules/no_unmarshalable_type"
	"github.com/strowk/vint/rules/no_unnecessary_blank_identifier"
	"github.com/strowk/vint/rules/no_unnecessary_block"
	"github.com/strowk/vint/rules/no_unnecessary_conversion"
	"github.com/strowk/vint/rules/no_unnecessary_defer"
	"github.com/strowk/vint/rules/no_unnecessary_defer_lambda"
	"github.com/strowk/vint/rules/no_unnecessary_deref"
	"github.com/strowk/vint/rules/no_unnecessary_format"
	"github.com/strowk/vint/rules/no_unnecessary_if"
	"github.com/strowk/vint/rules/no_unnecessary_lambda"
	"github.com/strowk/vint/rules/no_unnecessary_loop_var_copy"
	"github.com/strowk/vint/rules/no_unnecessary_stmt"
	"github.com/strowk/vint/rules/no_unnecessary_type_parens"
	"github.com/strowk/vint/rules/no_unobserved_field_assign"
	"github.com/strowk/vint/rules/no_unreachable_code"
	"github.com/strowk/vint/rules/no_unreachable_type_case"
	"github.com/strowk/vint/rules/no_unrecognized_directive"
	"github.com/strowk/vint/rules/no_unsafe_cors_bypass"
	"github.com/strowk/vint/rules/no_unsafe_deserialization"
	"github.com/strowk/vint/rules/no_unsafe_package"
	"github.com/strowk/vint/rules/no_unsafe_redirect_policy"
	"github.com/strowk/vint/rules/no_unsecure_url_scheme"
	"github.com/strowk/vint/rules/no_unsigned_negative_comparison"
	"github.com/strowk/vint/rules/no_untrappable_signal"
	"github.com/strowk/vint/rules/no_unused_constant"
	"github.com/strowk/vint/rules/no_unused_field"
	"github.com/strowk/vint/rules/no_unused_function"
	"github.com/strowk/vint/rules/no_unused_function_result"
	"github.com/strowk/vint/rules/no_unused_nolint"
	"github.com/strowk/vint/rules/no_unused_parameter"
	"github.com/strowk/vint/rules/no_unused_receiver"
	"github.com/strowk/vint/rules/no_unused_type"
	"github.com/strowk/vint/rules/no_unused_variable"
	"github.com/strowk/vint/rules/no_unused_write"
	"github.com/strowk/vint/rules/no_useless_break"
	"github.com/strowk/vint/rules/no_useless_fallthrough"
	"github.com/strowk/vint/rules/no_useless_math_call"
	"github.com/strowk/vint/rules/no_var_declaration"
	"github.com/strowk/vint/rules/no_variable_command_execution"
	"github.com/strowk/vint/rules/no_variable_shadowing"
	"github.com/strowk/vint/rules/no_wait_group_misuse"
	"github.com/strowk/vint/rules/no_weak_crypto_hash"
	"github.com/strowk/vint/rules/no_weak_encryption_algorithm"
	"github.com/strowk/vint/rules/no_weak_slice_guard"
	"github.com/strowk/vint/rules/no_writer_buffer_modification"
	"github.com/strowk/vint/rules/no_xss_taint"
	"github.com/strowk/vint/rules/no_yoda_condition"
	"github.com/strowk/vint/rules/no_zero_bytes_repeat"
	"github.com/strowk/vint/rules/no_zip_slip"
	"github.com/strowk/vint/rules/use_any"
	"github.com/strowk/vint/rules/use_assignment_operator"
	"github.com/strowk/vint/rules/use_buffer_string_or_bytes"
	"github.com/strowk/vint/rules/use_bytes_equal"
	"github.com/strowk/vint/rules/use_combined_append"
	"github.com/strowk/vint/rules/use_combined_param_type"
	"github.com/strowk/vint/rules/use_comment_spacing"
	"github.com/strowk/vint/rules/use_comment_spacings"
	"github.com/strowk/vint/rules/use_comments_density"
	"github.com/strowk/vint/rules/use_compiled_regexp"
	"github.com/strowk/vint/rules/use_consistent_receiver_name"
	"github.com/strowk/vint/rules/use_context_as_first_param"
	"github.com/strowk/vint/rules/use_convenience_func"
	"github.com/strowk/vint/rules/use_copy_builtin"
	"github.com/strowk/vint/rules/use_copy_for_slide"
	"github.com/strowk/vint/rules/use_decode_rune"
	"github.com/strowk/vint/rules/use_direct_method_call"
	"github.com/strowk/vint/rules/use_direct_return"
	"github.com/strowk/vint/rules/use_direct_string_comparison"
	"github.com/strowk/vint/rules/use_direct_string_range"
	"github.com/strowk/vint/rules/use_early_continue"
	"github.com/strowk/vint/rules/use_early_return"
	"github.com/strowk/vint/rules/use_else_if"
	"github.com/strowk/vint/rules/use_epoch_naming"
	"github.com/strowk/vint/rules/use_equal_fold"
	"github.com/strowk/vint/rules/use_error_last_return"
	"github.com/strowk/vint/rules/use_error_method"
	"github.com/strowk/vint/rules/use_error_naming"
	"github.com/strowk/vint/rules/use_errorf"
	"github.com/strowk/vint/rules/use_errors_as"
	"github.com/strowk/vint/rules/use_errors_new"
	"github.com/strowk/vint/rules/use_exported_comment"
	"github.com/strowk/vint/rules/use_file_header"
	"github.com/strowk/vint/rules/use_filename_format"
	"github.com/strowk/vint/rules/use_filepath_join"
	"github.com/strowk/vint/rules/use_fmt_print"
	"github.com/strowk/vint/rules/use_fprint"
	"github.com/strowk/vint/rules/use_func_doc_prefix"
	"github.com/strowk/vint/rules/use_getter_return"
	"github.com/strowk/vint/rules/use_http_no_body"
	"github.com/strowk/vint/rules/use_idiomatic_duration_name"
	"github.com/strowk/vint/rules/use_idiomatic_error_name"
	"github.com/strowk/vint/rules/use_idiomatic_naming"
	"github.com/strowk/vint/rules/use_idiomatic_receiver_name"
	"github.com/strowk/vint/rules/use_import_alias_naming"
	"github.com/strowk/vint/rules/use_indent_error_flow"
	"github.com/strowk/vint/rules/use_infinite_for"
	"github.com/strowk/vint/rules/use_inline_math_pow"
	"github.com/strowk/vint/rules/use_join_host_port"
	"github.com/strowk/vint/rules/use_loop_condition"
	"github.com/strowk/vint/rules/use_map_style"
	"github.com/strowk/vint/rules/use_matching_constant"
	"github.com/strowk/vint/rules/use_merged_conditional_decl"
	"github.com/strowk/vint/rules/use_merged_var_decl"
	"github.com/strowk/vint/rules/use_modern_octal_literal"
	"github.com/strowk/vint/rules/use_named_result"
	"github.com/strowk/vint/rules/use_optimal_field_alignment"
	"github.com/strowk/vint/rules/use_optimal_operands_order"
	"github.com/strowk/vint/rules/use_optimized_slice_clear"
	"github.com/strowk/vint/rules/use_package_comment"
	"github.com/strowk/vint/rules/use_package_comments"
	"github.com/strowk/vint/rules/use_package_naming"
	"github.com/strowk/vint/rules/use_parallel_assign_swap"
	"github.com/strowk/vint/rules/use_percent_q"
	"github.com/strowk/vint/rules/use_pointer_in_sync_pool"
	"github.com/strowk/vint/rules/use_printf_suffix"
	"github.com/strowk/vint/rules/use_raw_string_regexp"
	"github.com/strowk/vint/rules/use_receiver_naming"
	"github.com/strowk/vint/rules/use_regexp_must_compile"
	"github.com/strowk/vint/rules/use_repeated_arg_type_style"
	"github.com/strowk/vint/rules/use_short_var_decl"
	"github.com/strowk/vint/rules/use_simplified_bool_expr"
	"github.com/strowk/vint/rules/use_simplified_bool_return"
	"github.com/strowk/vint/rules/use_simplified_print_format"
	"github.com/strowk/vint/rules/use_simplified_regexp"
	"github.com/strowk/vint/rules/use_simplified_selector"
	"github.com/strowk/vint/rules/use_slice_append"
	"github.com/strowk/vint/rules/use_slice_style"
	"github.com/strowk/vint/rules/use_slices_sort"
	"github.com/strowk/vint/rules/use_standard_codegen_comment"
	"github.com/strowk/vint/rules/use_standard_deprecation_comment"
	"github.com/strowk/vint/rules/use_string_conversion_in_print"
	"github.com/strowk/vint/rules/use_string_format"
	"github.com/strowk/vint/rules/use_string_map_key"
	"github.com/strowk/vint/rules/use_string_writer"
	"github.com/strowk/vint/rules/use_strings_contains"
	"github.com/strowk/vint/rules/use_switch"
	"github.com/strowk/vint/rules/use_switch_default_style"
	"github.com/strowk/vint/rules/use_sync_map_load_and_delete"
	"github.com/strowk/vint/rules/use_tagged_switch"
	"github.com/strowk/vint/rules/use_time_date"
	"github.com/strowk/vint/rules/use_time_equal"
	"github.com/strowk/vint/rules/use_time_method"
	"github.com/strowk/vint/rules/use_time_naming"
	"github.com/strowk/vint/rules/use_time_since"
	"github.com/strowk/vint/rules/use_time_sleep"
	"github.com/strowk/vint/rules/use_time_until"
	"github.com/strowk/vint/rules/use_trim_function"
	"github.com/strowk/vint/rules/use_type_assert_result"
	"github.com/strowk/vint/rules/use_type_conversion"
	"github.com/strowk/vint/rules/use_type_def_first"
	"github.com/strowk/vint/rules/use_type_doc_prefix"
	"github.com/strowk/vint/rules/use_type_switch_chain"
	"github.com/strowk/vint/rules/use_type_switch_guard"
	"github.com/strowk/vint/rules/use_var_const_doc_prefix"
	"github.com/strowk/vint/rules/use_var_naming"
	"github.com/strowk/vint/rules/use_waitgroup_go"
	"github.com/strowk/vint/rules/use_write_byte"
)

var defaultRules = []lint.Rule{
	&no_var_declaration.VarDeclarationsRule{},
	&use_package_comments.PackageCommentsRule{},
	&no_blank_import.NoBlankImportRule{},
	&use_exported_comment.ExportedRule{},
	&use_var_naming.VarNamingRule{},
	&use_indent_error_flow.IndentErrorFlowRule{},
	&no_redundant_range_val.RangeRule{},
	&use_errorf.ErrorfRule{},
	&use_error_naming.ErrorNamingRule{},
	&no_error_strings.ErrorStringsRule{},
	&use_receiver_naming.ReceiverNamingRule{},
	&rule.IncrementDecrementRule{},
	&no_unexported_return.UnexportedReturnRule{},
	&use_time_naming.TimeNamingRule{},
	&no_context_keys_type.ContextKeysType{},
	&use_context_as_first_param.ContextAsArgumentRule{},
	&rule.EmptyBlockRule{},
	&no_superfluous_else.SuperfluousElseRule{},
	&no_unclosed_bodies.NoUnclosedBodiesRule{},
	&no_duplicate_code.NoDuplicateCodeRule{},
	&no_excessive_statements.NoExcessiveStatementsRule{},
	&no_high_cyclomatic_complexity.NoHighCyclomaticComplexityRule{},
	&no_long_functions.NoLongFunctionsRule{},
	&no_denied_import.NoDeniedImportRule{},
	&no_direct_error_comparison.NoDirectErrorComparisonRule{},
	&no_dynamic_errors.NoDynamicErrorsRule{},
	&no_file_scoped_denied_import.NoFileScopedDeniedImportRule{},
	&no_space_in_directive.NoSpaceInDirectiveRule{},
	&no_unallowed_import.NoUnallowedImportRule{},
	&no_unrecognized_directive.NoUnrecognizedDirectiveRule{},
	&no_unchecked_error.NoUncheckedErrorRule{},
	&no_range_variable_alias.NoRangeVariableAliasRule{},
	&no_slice_bounds_out_of_range.NoSliceBoundsOutOfRangeRule{},
	&no_atoi_overflow.NoAtoiOverflowRule{},
	&no_bind_to_all_interfaces.NoBindToAllInterfacesRule{},
	&no_exposed_pprof.NoExposedPprofRule{},
	&no_filesystem_root_serving.NoFilesystemRootServingRule{},
	&no_hardcoded_credentials.NoHardcodedCredentialsRule{},
	&no_hardcoded_iv.NoHardcodedIvRule{},
	&no_http_request_smuggling.NoHttpRequestSmugglingRule{},
	&no_insecure_cookie.NoInsecureCookieRule{},
	&no_insecure_host_key_callback.NoInsecureHostKeyCallbackRule{},
	&no_integer_overflow_conversion.NoIntegerOverflowConversionRule{},
	&no_missing_read_header_timeout.NoMissingReadHeaderTimeoutRule{},
	&no_serve_without_timeout.NoServeWithoutTimeoutRule{},
	&no_ssrf_via_variable.NoSsrfViaVariableRule{},
	&no_trojan_source_bidi.NoTrojanSourceBidiRule{},
	&no_unbounded_decompression.NoUnboundedDecompressionRule{},
	&no_unsafe_package.NoUnsafePackageRule{},
	&no_cgi_import.NoCgiImportRule{},
	&no_command_injection_taint.NoCommandInjectionTaintRule{},
	&no_context_propagation_failure.NoContextPropagationFailureRule{},
	&no_deprecated_hash_function.NoDeprecatedHashFunctionRule{},
	&no_deprecated_usage.NoDeprecatedUsageRule{},
	&no_file_inclusion_via_variable.NoFileInclusionViaVariableRule{},
	&no_filesystem_toctou.NoFilesystemToctouRule{},
	&no_insecure_random.NoInsecureRandomRule{},
	&no_insecure_tls_config.NoInsecureTlsConfigRule{},
	&no_log_injection_taint.NoLogInjectionTaintRule{},
	&no_path_traversal_taint.NoPathTraversalTaintRule{},
	&no_permissive_directory_permissions.NoPermissiveDirectoryPermissionsRule{},
	&no_permissive_file_permissions.NoPermissiveFilePermissionsRule{},
	&no_permissive_os_create.NoPermissiveOsCreateRule{},
	&no_permissive_write_file_permissions.NoPermissiveWriteFilePermissionsRule{},
	&no_predictable_temp_file.NoPredictableTempFileRule{},
	&no_secret_in_serialization.NoSecretInSerializationRule{},
	&no_short_rsa_key.NoShortRsaKeyRule{},
	&no_smtp_injection_taint.NoSmtpInjectionTaintRule{},
	&no_sql_concatenation.NoSqlConcatenationRule{},
	&no_sql_format_string.NoSqlFormatStringRule{},
	&no_sql_injection_taint.NoSqlInjectionTaintRule{},
	&no_ssh_auth_bypass.NoSshAuthBypassRule{},
	&no_ssrf_taint.NoSsrfTaintRule{},
	&no_template_injection.NoTemplateInjectionRule{},
	&no_tls_conn_handshake.NoTlsConnHandshakeRule{},
	&no_tls_dial.NoTlsDialRule{},
	&no_tls_dial_with_dialer.NoTlsDialWithDialerRule{},
	&no_tls_session_resumption_bypass.NoTlsSessionResumptionBypassRule{},
	&no_unbounded_form_parsing.NoUnboundedFormParsingRule{},
	&no_unsafe_cors_bypass.NoUnsafeCorsBypassRule{},
	&no_unsafe_redirect_policy.NoUnsafeRedirectPolicyRule{},
	&no_unescaped_html_template.NoUnescapedHtmlTemplateRule{},
	&no_unsafe_deserialization.NoUnsafeDeserializationRule{},
	&no_variable_command_execution.NoVariableCommandExecutionRule{},
	&no_weak_crypto_hash.NoWeakCryptoHashRule{},
	&no_weak_encryption_algorithm.NoWeakEncryptionAlgorithmRule{},
	&no_xss_taint.NoXssTaintRule{},
	&no_zip_slip.NoZipSlipRule{},
	&no_excessive_blank_identifiers.NoExcessiveBlankIdentifiersRule{},
	&no_init_function.NoInitFunctionRule{},
	&no_dot_import.NoDotImportRule{},
	&no_repeated_strings.NoRepeatedStringsRule{},
	&no_unnecessary_conversion.NoUnnecessaryConversionRule{},
	&no_unnecessary_deref.NoUnnecessaryDerefRule{},
	&no_unnecessary_blank_identifier.NoUnnecessaryBlankIdentifierRule{},
	&no_unnecessary_loop_var_copy.NoUnnecessaryLoopVarCopyRule{},
	&use_printf_suffix.UsePrintfSuffixRule{},
	&no_single_arg_append.NoSingleArgAppendRule{},
	&no_asm_decl_mismatch.NoAsmDeclMismatchRule{},
	&no_self_assignment.NoSelfAssignmentRule{},
	&no_self_referencing_finalizer.NoSelfReferencingFinalizerRule{},
	&no_atomic_assign_misuse.NoAtomicAssignMisuseRule{},
	&no_malformed_build_tag.NoMalformedBuildTagRule{},
	&no_malformed_directive.NoMalformedDirectiveRule{},
	&no_malformed_struct_tag.NoMalformedStructTagRule{},
	&no_cgo_pointer_violation.NoCgoPointerViolationRule{},
	&no_unkeyed_literal.NoUnkeyedLiteralRule{},
	&no_copied_lock.NoCopiedLockRule{},
	&no_defer_time_misuse.NoDeferTimeMisuseRule{},
	&no_incorrect_time_format.NoIncorrectTimeFormatRule{},
	&no_http_client_get.NoHttpClientGetRule{},
	&no_http_client_head.NoHttpClientHeadRule{},
	&no_http_client_post.NoHttpClientPostRule{},
	&no_http_client_post_form.NoHttpClientPostFormRule{},
	&no_http_get.NoHttpGetRule{},
	&no_http_head.NoHttpHeadRule{},
	&no_http_post.NoHttpPostRule{},
	&no_http_post_form.NoHttpPostFormRule{},
	&no_http_new_request.NoHttpNewRequestRule{},
	&no_http_response_misuse.NoHttpResponseMisuseRule{},
	&no_httptest_new_request.NoHttptestNewRequestRule{},
	&no_invalid_binary_arg.NoInvalidBinaryArgRule{},
	&no_invalid_errors_as.NoInvalidErrorsAsRule{},
	&no_frame_pointer_clobber.NoFramePointerClobberRule{},
	&use_join_host_port.UseJoinHostPortRule{},
	&no_impossible_interface_assert.NoImpossibleInterfaceAssertRule{},
	&no_loop_closure_capture.NoLoopClosureCaptureRule{},
	&no_lost_cancel.NoLostCancelRule{},
	&no_wait_group_misuse.NoWaitGroupMisuseRule{},
	&no_writer_buffer_modification.NoWriterBufferModificationRule{},
	&no_nil_context.NoNilContextRule{},
	&no_nil_func_comparison.NoNilFuncComparisonRule{},
	&no_nil_map_assignment.NoNilMapAssignmentRule{},
	&no_printf_format_mismatch.NoPrintfFormatMismatchRule{},
	&no_excessive_shift.NoExcessiveShiftRule{},
	&no_dubious_bit_shift.NoDubiousBitShiftRule{},
	&no_unbuffered_signal_channel.NoUnbufferedSignalChannelRule{},
	&no_untrappable_signal.NoUntrappableSignalRule{},
	&no_slog_key_value_mismatch.NoSlogKeyValueMismatchRule{},
	&no_odd_size_slice_arg.NoOddSizeSliceArgRule{},
	&no_std_method_signature_mismatch.NoStdMethodSignatureMismatchRule{},
	&no_stdlib_version_mismatch.NoStdlibVersionMismatchRule{},
	&no_string_int_conversion.NoStringIntConversionRule{},
	&no_test_fatal_in_goroutine.NoTestFatalInGoroutineRule{},
	&no_test_main_without_exit.NoTestMainWithoutExitRule{},
	&no_malformed_test_function.NoMalformedTestFunctionRule{},
	&no_non_canonical_header_key.NoNonCanonicalHeaderKeyRule{},
	&no_redundant_canonical_header_key.NoRedundantCanonicalHeaderKeyRule{},
	&no_non_pointer_unmarshal.NoNonPointerUnmarshalRule{},
	&no_unmarshalable_struct.NoUnmarshalableStructRule{},
	&no_unmarshalable_type.NoUnmarshalableTypeRule{},
	&no_unreachable_code.NoUnreachableCodeRule{},
	&no_invalid_unsafe_pointer.NoInvalidUnsafePointerRule{},
	&no_unused_function_result.NoUnusedFunctionResultRule{},
	&no_unused_write.NoUnusedWriteRule{},
	&no_ineffectual_assignment.NoIneffectualAssignmentRule{},
	&no_ineffective_bitwise_op.NoIneffectiveBitwiseOpRule{},
	&no_ineffective_break.NoIneffectiveBreakRule{},
	&no_ineffective_rand_call.NoIneffectiveRandCallRule{},
	&no_line_too_long.NoLineTooLongRule{},
	&no_magic_number_in_argument.NoMagicNumberInArgumentRule{},
	&no_magic_number_in_assignment.NoMagicNumberInAssignmentRule{},
	&no_magic_number_in_case.NoMagicNumberInCaseRule{},
	&no_magic_number_in_condition.NoMagicNumberInConditionRule{},
	&no_magic_number_in_operation.NoMagicNumberInOperationRule{},
	&no_magic_number_in_return.NoMagicNumberInReturnRule{},

	// temporary disable while researching optimization, this one is weirdly slow
	&no_misspelled_words.NoMisspelledWordsRule{},

	&no_nan_comparison.NoNanComparisonRule{},
	&no_naked_return.NoNakedReturnRule{},
	&no_net_dial.NoNetDialRule{},
	&no_net_dial_timeout.NoNetDialTimeoutRule{},
	&no_net_ip_bytes_equal.NoNetIpBytesEqualRule{},
	&no_net_listen.NoNetListenRule{},
	&no_net_listen_packet.NoNetListenPacketRule{},
	&no_net_lookup_cname.NoNetLookupCnameRule{},
	&no_net_lookup_host.NoNetLookupHostRule{},
	&no_net_lookup_ip.NoNetLookupIpRule{},
	&no_net_lookup_port.NoNetLookupPortRule{},
	&no_net_lookup_mx.NoNetLookupMxRule{},
	&no_net_lookup_srv.NoNetLookupSrvRule{},
	&no_net_lookup_ns.NoNetLookupNsRule{},
	&no_net_lookup_txt.NoNetLookupTxtRule{},
	&no_net_lookup_addr.NoNetLookupAddrRule{},
	&no_sql_db_begin.NoSqlDbBeginRule{},
	&no_sql_db_exec.NoSqlDbExecRule{},
	&no_sql_db_ping.NoSqlDbPingRule{},
	&no_sql_db_prepare.NoSqlDbPrepareRule{},
	&no_sql_db_query.NoSqlDbQueryRule{},
	&no_sql_db_query_row.NoSqlDbQueryRowRule{},
	&no_sql_tx_exec.NoSqlTxExecRule{},
	&no_sql_tx_prepare.NoSqlTxPrepareRule{},
	&no_sql_tx_query.NoSqlTxQueryRule{},
	&no_sql_tx_query_row.NoSqlTxQueryRowRule{},
	&no_sql_tx_stmt.NoSqlTxStmtRule{},
	&no_sql_stmt_exec.NoSqlStmtExecRule{},
	&no_sql_stmt_query.NoSqlStmtQueryRule{},
	&no_sql_stmt_query_row.NoSqlStmtQueryRowRule{},
	&no_exec_command.NoExecCommandRule{},
	&no_invalid_exec_command_arg.NoInvalidExecCommandArgRule{},
	&no_nolint_non_machine_readable.NoNolintNonMachineReadableRule{},
	&no_nolint_parse_error.NoNolintParseErrorRule{},
	&no_nolint_without_specific_linter.NoNolintWithoutSpecificLinterRule{},
	&no_nolint_without_explanation.NoNolintWithoutExplanationRule{},
	&no_unused_nolint.NoUnusedNolintRule{},
	&no_unused_parameter.NoUnusedParameterRule{},
	&no_constant_parameter.NoConstantParameterRule{},
	&no_constant_result.NoConstantResultRule{},
	&no_unused_function.NoUnusedFunctionRule{},
	&no_unused_type.NoUnusedTypeRule{},
	&no_unused_variable.NoUnusedVariableRule{},
	&no_unused_constant.NoUnusedConstantRule{},

	// temporary off for performance research
	&no_unused_field.NoUnusedFieldRule{},

	&no_leading_blank_line.NoLeadingBlankLineRule{},
	&no_trailing_blank_line.NoTrailingBlankLineRule{},
	&no_discarded_append.NoDiscardedAppendRule{},
	&no_mismatched_append_assign.NoMismatchedAppendAssignRule{},
	&use_combined_append.UseCombinedAppendRule{},
	&no_swapped_arguments.NoSwappedArgumentsRule{},
	&use_assignment_operator.UseAssignmentOperatorRule{},
	&no_noop_function_call.NoNoopFunctionCallRule{},
	&no_impossible_condition.NoImpossibleConditionRule{},
	&no_impossible_builtin_result.NoImpossibleBuiltinResultRule{},
	&no_capitalized_local.NoCapitalizedLocalRule{},
	&no_capitalized_error_string.NoCapitalizedErrorStringRule{},
	&no_unreachable_type_case.NoUnreachableTypeCaseRule{},
	&no_misplaced_default_case.NoMisplacedDefaultCaseRule{},
	&use_standard_codegen_comment.UseStandardCodegenCommentRule{},
	&use_standard_deprecation_comment.UseStandardDeprecationCommentRule{},
	&use_comment_spacing.UseCommentSpacingRule{},
	&no_duplicate_build_constraint.NoDuplicateBuildConstraintRule{},
	&no_duplicate_argument.NoDuplicateArgumentRule{},
	&no_duplicate_branch_body.NoDuplicateBranchBodyRule{},
	&no_duplicate_case.NoDuplicateCaseRule{},
	&no_duplicate_sub_expression.NoDuplicateSubExpressionRule{},
	&no_duplicate_cutset_chars.NoDuplicateCutsetCharsRule{},
	&no_duplicate_if_condition.NoDuplicateIfConditionRule{},
	&use_else_if.UseElseIfRule{},
	&no_defer_close_before_err_check.NoDeferCloseBeforeErrCheckRule{},
	&no_defer_in_infinite_loop.NoDeferInInfiniteLoopRule{},
	&no_exit_after_defer.NoExitAfterDeferRule{},
	&no_flag_deref_before_parse.NoFlagDerefBeforeParseRule{},
	&no_invalid_flag_name.NoInvalidFlagNameRule{},
	&no_huge_param.NoHugeParamRule{},
	&no_string_index_allocation.NoStringIndexAllocationRule{},
	&use_string_map_key.UseStringMapKeyRule{},
	&no_redundant_string_byte_conversion.NoRedundantStringByteConversionRule{},
	&use_inline_math_pow.UseInlineMathPowRule{},
	&use_switch.UseSwitchRule{},
	&use_tagged_switch.UseTaggedSwitchRule{},
	&no_suspicious_map_key.NoSuspiciousMapKeyRule{},
	&no_suspicious_time_sleep.NoSuspiciousTimeSleepRule{},
	&no_temp_dir_deletion.NoTempDirDeletionRule{},
	&no_empty_branch.NoEmptyBranchRule{},
	&no_empty_critical_section.NoEmptyCriticalSectionRule{},
	&no_time_tick.NoTimeTickRule{},
	&no_timer_reset_retval.NoTimerResetRetvalRule{},
	&no_immediate_new_deref.NoImmediateNewDerefRule{},
	&no_off_by_one_error.NoOffByOneErrorRule{},
	&no_range_expr_copy.NoRangeExprCopyRule{},
	&no_range_val_copy.NoRangeValCopyRule{},
	&no_redundant_rune_conversion.NoRedundantRuneConversionRule{},
	&use_regexp_must_compile.UseRegexpMustCompileRule{},
	&use_compiled_regexp.UseCompiledRegexpRule{},
	&no_single_case_select.NoSingleCaseSelectRule{},
	&no_single_case_switch.NoSingleCaseSwitchRule{},
	&no_single_iteration_loop.NoSingleIterationLoopRule{},
	&no_invariant_loop_condition.NoInvariantLoopConditionRule{},
	&no_redundant_switch_true.NoRedundantSwitchTrueRule{},
	&no_always_true_len_check.NoAlwaysTrueLenCheckRule{},
	&no_redundant_type_assertion.NoRedundantTypeAssertionRule{},
	&use_type_switch_guard.UseTypeSwitchGuardRule{},
	&use_type_assert_result.UseTypeAssertResultRule{},
	&no_unnecessary_lambda.NoUnnecessaryLambdaRule{},
	&no_redundant_slice_expression.NoRedundantSliceExpressionRule{},
	&use_parallel_assign_swap.UseParallelAssignSwapRule{},
	&use_convenience_func.UseConvenienceFuncRule{},
	&no_zero_bytes_repeat.NoZeroBytesRepeatRule{},
	&no_invalid_host_port.NoInvalidHostPortRule{},
	&no_invalid_regexp.NoInvalidRegexpRule{},
	&no_invalid_template.NoInvalidTemplateRule{},
	&no_invalid_url_parse.NoInvalidUrlParseRule{},
	&no_invalid_utf8_string_arg.NoInvalidUtf8StringArgRule{},
	&no_inappropriate_context_key.NoInappropriateContextKeyRule{},
	&no_invalid_strconv_arg.NoInvalidStrconvArgRule{},
	&no_overlapping_encoder_slice.NoOverlappingEncoderSliceRule{},
	&no_benchmark_n_assignment.NoBenchmarkNAssignmentRule{},
	&no_address_of_dereference.NoAddressOfDereferenceRule{},
	&no_address_nil_comparison.NoAddressNilComparisonRule{},
	&no_unsigned_negative_comparison.NoUnsignedNegativeComparisonRule{},
	&no_unobserved_field_assign.NoUnobservedFieldAssignRule{},
	&no_overwritten_argument.NoOverwrittenArgumentRule{},
	&no_useless_math_call.NoUselessMathCallRule{},
	&no_impossible_nil_comparison.NoImpossibleNilComparisonRule{},
	&no_integer_division_truncation.NoIntegerDivisionTruncationRule{},
	&no_imprecise_constant.NoImpreciseConstantRule{},
	&no_implicit_const_value.NoImplicitConstValueRule{},
	&no_ignored_query_modification.NoIgnoredQueryModificationRule{},
	&no_modulo_one.NoModuloOneRule{},
	&no_never_nil_check.NoNeverNilCheckRule{},
	&no_non_octal_file_mode.NoNonOctalFileModeRule{},
	&no_empty_for_loop.NoEmptyForLoopRule{},
	&no_select_break_confusion.NoSelectBreakConfusionRule{},
	&no_infinite_recursion.NoInfiniteRecursionRule{},
	&use_pointer_in_sync_pool.UsePointerInSyncPoolRule{},
	&no_type_assert_else_misread.NoTypeAssertElseMisreadRule{},
	&use_copy_builtin.UseCopyBuiltinRule{},
	&use_copy_for_slide.UseCopyForSlideRule{},
	&no_explicit_bool_comparison.NoExplicitBoolComparisonRule{},
	&use_strings_contains.UseStringsContainsRule{},
	&use_bytes_equal.UseBytesEqualRule{},
	&use_infinite_for.UseInfiniteForRule{},
	&use_raw_string_regexp.UseRawStringRegexpRule{},
	&use_simplified_bool_return.UseSimplifiedBoolReturnRule{},
	&no_redundant_make_args.NoRedundantMakeArgsRule{},
	&no_redundant_nil_loop_check.NoRedundantNilLoopCheckRule{},
	&no_redundant_nil_slice_check.NoRedundantNilSliceCheckRule{},
	&no_redundant_nil_type_check.NoRedundantNilTypeCheckRule{},
	&no_default_slice_index.NoDefaultSliceIndexRule{},
	&use_slice_append.UseSliceAppendRule{},
	&use_time_equal.UseTimeEqualRule{},
	&use_time_since.UseTimeSinceRule{},
	&use_time_sleep.UseTimeSleepRule{},
	&use_time_until.UseTimeUntilRule{},
	&use_type_conversion.UseTypeConversionRule{},
	&use_trim_function.UseTrimFunctionRule{},
	&use_merged_conditional_decl.UseMergedConditionalDeclRule{},
	&use_merged_var_decl.UseMergedVarDeclRule{},
	&no_redundant_control_flow.NoRedundantControlFlowRule{},
	&use_error_method.UseErrorMethodRule{},
	&use_error_last_return.UseErrorLastReturnRule{},
	&use_direct_string_range.UseDirectStringRangeRule{},
	&use_buffer_string_or_bytes.UseBufferStringOrBytesRule{},
	&no_guard_around_delete.NoGuardAroundDeleteRule{},
	&no_guard_around_map_access.NoGuardAroundMapAccessRule{},
	&use_simplified_print_format.UseSimplifiedPrintFormatRule{},
	&use_string_conversion_in_print.UseStringConversionInPrintRule{},
	&use_package_comment.UsePackageCommentRule{},
	&use_idiomatic_error_name.UseIdiomaticErrorNameRule{},
	&use_idiomatic_naming.UseIdiomaticNamingRule{},
	&use_idiomatic_receiver_name.UseIdiomaticReceiverNameRule{},
	&use_idiomatic_duration_name.UseIdiomaticDurationNameRule{},
	&use_consistent_receiver_name.UseConsistentReceiverNameRule{},
	&no_control_char_in_string.NoControlCharInStringRule{},
	&use_func_doc_prefix.UseFuncDocPrefixRule{},
	&use_type_doc_prefix.UseTypeDocPrefixRule{},
	&use_var_const_doc_prefix.UseVarConstDocPrefixRule{},
	&no_redundant_var_type.NoRedundantVarTypeRule{},
	&use_loop_condition.UseLoopConditionRule{},
	&use_simplified_selector.UseSimplifiedSelectorRule{},
	&no_non_wrapping_errorf.NoNonWrappingErrorfRule{},
	&use_errors_as.UseErrorsAsRule{},
}

var allRules = append([]lint.Rule{
	&no_excessive_arguments.NoExcessiveArgumentsRule{},
	&use_file_header.FileHeaderRule{},
	&no_confusing_naming.ConfusingNamingRule{},
	&use_getter_return.GetReturnRule{},
	&no_modified_parameter.ModifiesParamRule{},
	&no_confusing_results.ConfusingResultsRule{},
	&no_deep_exit.DeepExitRule{},
	&no_flag_parameter.FlagParamRule{},
	&no_unnecessary_stmt.UnnecessaryStmtRule{},
	&rule.StructTagRule{},
	&no_modified_value_receiver.ModifiesValRecRule{},
	&rule.ConstantLogicalExprRule{},
	&no_excessive_function_results.FunctionResultsLimitRule{},
	&no_excessive_public_structs.NoExcessivePublicStructsRule{},
	&no_call_to_gc.CallToGCRule{},
	&no_duplicated_imports.DuplicatedImportsRule{},
	&rule.ImportShadowingRule{},
	&rule.BareReturnRule{},
	&no_unused_receiver.UnusedReceiverRule{},
	&no_high_cognitive_complexity.NoHighCognitiveComplexityRule{},
	&use_string_format.StringFormatRule{},
	&use_early_return.EarlyReturnRule{},
	&no_unconditional_recursion.UnconditionalRecursionRule{},
	&no_identical_branches.IdenticalBranchesRule{},
	&no_defer_gotcha.DeferRule{},
	&no_unexported_naming.UnexportedNamingRule{},
	&no_nested_structs.NestedStructs{},
	&no_useless_break.UselessBreak{},
	&use_time_date.TimeDateRule{},
	&no_banned_characters.NoBannedCharactersRule{},
	&use_optimal_operands_order.OptimizeOperandsOrderRule{},
	&use_any.UseAnyRule{},
	&rule.DataRaceRule{},
	&use_comment_spacings.CommentSpacingsRule{},
	&use_direct_return.IfReturnRule{},
	&no_redundant_import_alias.NoRedundantImportAliasRule{},
	&use_import_alias_naming.ImportAliasNamingRule{},
	&use_map_style.EnforceMapStyleRule{},
	&use_repeated_arg_type_style.EnforceRepeatedArgTypeStyleRule{},
	&use_slice_style.EnforceSliceStyleRule{},
	&no_excessive_control_nesting.MaxControlNestingRule{},
	&use_comments_density.CommentsDensityRule{},
	&no_excessive_file_length.NoExcessiveFileLengthRule{},
	&use_filename_format.FilenameFormatRule{},
	&no_redundant_build_tag.RedundantBuildTagRule{},
	&use_errors_new.UseErrorsNewRule{},
	&no_redundant_test_main_exit.RedundantTestMainExitRule{},
	&no_unnecessary_format.UnnecessaryFormatRule{},
	&use_fmt_print.UseFmtPrintRule{},
	&use_switch_default_style.EnforceSwitchStyleRule{},
	&no_identical_switch_conditions.IdenticalSwitchConditionsRule{},
	&no_identical_if_else_if_conditions.IdenticalIfElseIfConditionsRule{},
	&no_identical_if_else_if_branches.IdenticalIfElseIfBranchesRule{},
	&no_identical_switch_branches.IdenticalSwitchBranchesRule{},
	&no_useless_fallthrough.UselessFallthroughRule{},
	&no_package_directory_mismatch.PackageDirectoryMismatchRule{},
	&use_waitgroup_go.UseWaitGroupGoRule{},
	&no_unsecure_url_scheme.UnsecureURLSchemeRule{},
	&no_inefficient_map_lookup.InefficientMapLookupRule{},
	&no_forbidden_call_in_wg_go.ForbiddenCallInWgGoRule{},
	&no_unnecessary_if.UnnecessaryIfRule{},
	&use_epoch_naming.EpochNamingRule{},
	&use_slices_sort.UseSlicesSort{},
	&use_package_naming.PackageNamingRule{},
	&no_atomic_alignment_issue.NoAtomicAlignmentIssueRule{},
	&no_blank_error_assignment.NoBlankErrorAssignmentRule{},
	&no_specific_function_call.NoSpecificFunctionCallRule{},
	&no_deep_equal_errors.NoDeepEqualErrorsRule{},
	&no_duplicate_constants.NoDuplicateConstantsRule{},
	&no_unchecked_type_assertion.NoUncheckedTypeAssertionRule{},
	&no_repeated_numbers.NoRepeatedNumbersRule{},
	&use_matching_constant.UseMatchingConstantRule{},
	&use_optimal_field_alignment.UseOptimalFieldAlignmentRule{},
	&use_optimized_slice_clear.UseOptimizedSliceClearRule{},
	&no_conflicting_http_mux_patterns.NoConflictingHttpMuxPatternsRule{},
	&no_invalid_sort_slice_arg.NoInvalidSortSliceArgRule{},
	&no_nil_dereference.NoNilDereferenceRule{},
	&no_nil_variable_return.NoNilVariableReturnRule{},
	&no_variable_shadowing.NoVariableShadowingRule{},
	&no_multi_line_if_break.NoMultiLineIfBreakRule{},
	&no_multi_line_func_break.NoMultiLineFuncBreakRule{},
	&no_bad_lock_pattern.NoBadLockPatternRule{},
	&no_bad_regexp_pattern.NoBadRegexpPatternRule{},
	&no_bad_sort_usage.NoBadSortUsageRule{},
	&no_inline_sync_once_func.NoInlineSyncOnceFuncRule{},
	&use_simplified_bool_expr.UseSimplifiedBoolExprRule{},
	&no_builtin_shadow.NoBuiltinShadowRule{},
	&no_builtin_shadow_decl.NoBuiltinShadowDeclRule{},
	&no_import_shadow.NoImportShadowRule{},
	&use_direct_method_call.UseDirectMethodCallRule{},
	&use_direct_string_comparison.UseDirectStringComparisonRule{},
	&no_commented_out_code.NoCommentedOutCodeRule{},
	&no_commented_out_import.NoCommentedOutImportRule{},
	&no_defer_in_loop.NoDeferInLoopRule{},
	&no_dynamic_format_string.NoDynamicFormatStringRule{},
	&no_unnecessary_block.NoUnnecessaryBlockRule{},
	&no_unnecessary_defer.NoUnnecessaryDeferRule{},
	&no_unnecessary_defer_lambda.NoUnnecessaryDeferLambdaRule{},
	&no_unnecessary_type_parens.NoUnnecessaryTypeParensRule{},
	&no_doc_comment_stub.NoDocCommentStubRule{},
	&no_duplicate_import.NoDuplicateImportRule{},
	&no_duplicate_option.NoDuplicateOptionRule{},
	&no_empty_declaration.NoEmptyDeclarationRule{},
	&no_empty_fallthrough.NoEmptyFallthroughRule{},
	&no_eval_order_dependency.NoEvalOrderDependencyRule{},
	&no_exposed_sync_mutex.NoExposedSyncMutexRule{},
	&no_external_error_reassign.NoExternalErrorReassignRule{},
	&use_equal_fold.UseEqualFoldRule{},
	&no_separator_in_filepath_join.NoSeparatorInFilepathJoinRule{},
	&use_filepath_join.UseFilepathJoinRule{},
	&no_mixed_case_hex_literal.NoMixedCaseHexLiteralRule{},
	&use_http_no_body.UseHttpNoBodyRule{},
	&no_side_effect_in_init_clause.NoSideEffectInInitClauseRule{},
	&use_early_continue.UseEarlyContinueRule{},
	&use_modern_octal_literal.UseModernOctalLiteralRule{},
	&no_pointer_to_ref_param.NoPointerToRefParamRule{},
	&no_range_append_all.NoRangeAppendAllRule{},
	&no_redundant_label.NoRedundantLabelRule{},
	&no_redundant_sprint.NoRedundantSprintRule{},
	&no_redundant_string_concat.NoRedundantStringConcatRule{},
	&use_combined_param_type.UseCombinedParamTypeRule{},
	&use_decode_rune.UseDecodeRuneRule{},
	&use_fprint.UseFprintRule{},
	&use_string_writer.UseStringWriterRule{},
	&use_write_byte.UseWriteByteRule{},
	&no_unescaped_regexp_dot.NoUnescapedRegexpDotRule{},
	&use_simplified_regexp.UseSimplifiedRegexpRule{},
	&no_ignored_query_result.NoIgnoredQueryResultRule{},
	&no_missing_return_after_http_error.NoMissingReturnAfterHttpErrorRule{},
	&no_suspicious_sort_slice.NoSuspiciousSortSliceRule{},
	&use_short_var_decl.UseShortVarDeclRule{},
	&use_percent_q.UsePercentQRule{},
	&no_strings_compare.NoStringsCompareRule{},
	&use_sync_map_load_and_delete.UseSyncMapLoadAndDeleteRule{},
	&use_time_method.UseTimeMethodRule{},
	&no_todo_without_detail.NoTodoWithoutDetailRule{},
	&no_excessive_results.NoExcessiveResultsRule{},
	&no_truncating_comparison.NoTruncatingComparisonRule{},
	&no_unchecked_inline_error.NoUncheckedInlineErrorRule{},
	&no_weak_slice_guard.NoWeakSliceGuardRule{},
	&no_yoda_condition.NoYodaConditionRule{},
	&use_type_def_first.UseTypeDefFirstRule{},
	&use_type_switch_chain.UseTypeSwitchChainRule{},
	&use_named_result.UseNamedResultRule{},
}, defaultRules...)

func GetAllRules() []lint.Rule {
	return allRules
}

// allFormatters is a list of all available formatters to output the linting results.
// Keep the list sorted and in sync with available formatters in README.md.
var allFormatters = []lint.Formatter{
	&formatter.Checkstyle{},
	&formatter.Default{},
	&formatter.Friendly{},
	&formatter.JSON{},
	&formatter.NDJSON{},
	&formatter.Plain{},
	&formatter.Sarif{},
	&formatter.Stylish{},
	&formatter.Unix{},
}

func getFormatters() map[string]lint.Formatter {
	result := map[string]lint.Formatter{}
	for _, f := range allFormatters {
		result[f.Name()] = f
	}
	return result
}

// GetLintingRules yields the linting rules that must be applied by the linter.
func GetLintingRules(config *lint.Config, extraRules []lint.Rule) ([]lint.Rule, error) {
	rulesMap := map[string]lint.Rule{}
	for _, r := range allRules {
		rulesMap[lint.FullRuleName(r)] = r
	}
	for _, r := range extraRules {
		fullName := lint.FullRuleName(r)
		if _, ok := rulesMap[fullName]; ok {
			continue
		}
		rulesMap[fullName] = r
	}

	var lintingRules []lint.Rule
	for name, ruleConfig := range config.Rules {
		actualName := actualRuleName(name)
		r, ok := rulesMap[actualName]
		if !ok {
			return nil, fmt.Errorf("cannot find rule: %s", name)
		}

		if ruleConfig.Disabled {
			continue // skip disabled rules
		}

		if r, ok := r.(lint.ConfigurableRule); ok {
			if err := r.Configure(ruleConfig.Arguments); err != nil {
				return nil, fmt.Errorf("cannot configure rule: %q: %w", name, err)
			}
		}

		lintingRules = append(lintingRules, r)
	}

	return lintingRules, nil
}

func actualRuleName(name string) string {
	return name
}

func parseConfig(data []byte, config *lint.Config) error {
	err := toml.Unmarshal(data, config)
	if err != nil {
		return fmt.Errorf("cannot parse the config file: %w", err)
	}
	for k, r := range config.Rules {
		err := r.Initialize()
		if err != nil {
			return fmt.Errorf("error in config of rule [%s] : [%w]", k, err)
		}
		config.Rules[k] = r
	}

	return nil
}

// vintYAMLConfig represents the vint.yaml configuration file structure.
type vintYAMLConfig struct {
	Settings map[string]vintYAMLRuleConfig `yaml:"settings"`
}

// vintYAMLRuleConfig represents a single rule's configuration in vint.yaml.
type vintYAMLRuleConfig struct {
	Severity  string   `yaml:"severity"`
	Disabled  bool     `yaml:"disabled"`
	Arguments []any    `yaml:"arguments"`
	Exclude   []string `yaml:"exclude"`
}

// mergeVintYAML reads vint.yaml from the given path and merges its settings
// into the existing config. YAML settings override matching TOML settings.
func mergeVintYAML(path string, config *lint.Config) error {
	data, err := os.ReadFile(path) //nolint:gosec // ignore G304: potential file inclusion via variable
	if err != nil {
		return fmt.Errorf("cannot read vint.yaml: %w", err)
	}

	var yamlCfg vintYAMLConfig
	if err := yaml.Unmarshal(data, &yamlCfg); err != nil {
		return fmt.Errorf("cannot parse vint.yaml: %w", err)
	}

	if config.Rules == nil {
		config.Rules = map[string]lint.RuleConfig{}
	}

	for name, yamlRule := range yamlCfg.Settings {
		rc := lint.RuleConfig{
			Severity:  lint.Severity(yamlRule.Severity),
			Disabled:  yamlRule.Disabled,
			Arguments: yamlRule.Arguments,
			Exclude:   yamlRule.Exclude,
		}
		if err := rc.Initialize(); err != nil {
			return fmt.Errorf("error in vint.yaml config of rule [%s]: %w", name, err)
		}
		config.Rules[name] = rc
	}

	return nil
}

func validateConfig(config *lint.Config) error {
	if config.EnableAllRules && config.EnableDefaultRules {
		return errors.New("config options enableAllRules and enableDefaultRules cannot be combined")
	}
	return nil
}

func normalizeConfig(config *lint.Config) {
	if len(config.Rules) == 0 {
		config.Rules = map[string]lint.RuleConfig{}
	}

	addRules := func(config *lint.Config, rules []lint.Rule) {
		for _, r := range rules {
			ruleName := lint.FullRuleName(r)
			if _, ok := config.Rules[ruleName]; !ok {
				config.Rules[ruleName] = lint.RuleConfig{}
			}
		}
	}

	if config.EnableAllRules {
		addRules(config, allRules)
	} else if config.EnableDefaultRules {
		addRules(config, defaultRules)
	}

	severity := config.Severity
	if severity != "" {
		for k, v := range config.Rules {
			if v.Severity == "" {
				v.Severity = severity
			}
			config.Rules[k] = v
		}
		for k, v := range config.Directives {
			if v.Severity == "" {
				v.Severity = severity
			}
			config.Directives[k] = v
		}
	}
}

const defaultConfidence = 0.8

// vintYAMLPath is the hardcoded config file name for YAML configuration.
const vintYAMLPath = "vint.yaml"

// GetConfig yields the configuration.
func GetConfig(configPath string) (*lint.Config, error) {
	config := &lint.Config{}
	switch {
	case configPath != "":
		config.Confidence = defaultConfidence
		data, err := os.ReadFile(configPath) //nolint:gosec // ignore G304: potential file inclusion via variable
		if err != nil {
			return nil, errors.New("cannot read the config file")
		}
		err = parseConfig(data, config)
		if err != nil {
			return nil, err
		}

	default: // no configuration provided
		config = defaultConfig()
	}

	if err := validateConfig(config); err != nil {
		return nil, err
	}

	normalizeConfig(config)

	// Merge vint.yaml if it exists in CWD (overrides TOML settings)
	if _, err := os.Stat(vintYAMLPath); err == nil {
		if err := mergeVintYAML(vintYAMLPath, config); err != nil {
			return nil, err
		}
	}

	return config, nil
}

// GetFormatter yields the formatter for lint failures.
func GetFormatter(formatterName string) (lint.Formatter, error) {
	formatters := getFormatters()
	if formatterName == "" {
		return formatters["default"], nil
	}
	f, ok := formatters[formatterName]
	if !ok {
		return nil, fmt.Errorf("unknown formatter %v", formatterName)
	}
	return f, nil
}

func defaultConfig() *lint.Config {
	defaultConfig := lint.Config{
		Confidence: defaultConfidence,
		Severity:   lint.SeverityWarning,
		Rules:      map[string]lint.RuleConfig{},
	}
	for _, r := range defaultRules {
		defaultConfig.Rules[lint.FullRuleName(r)] = lint.RuleConfig{}
	}
	return &defaultConfig
}
