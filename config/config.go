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
	"github.com/strowk/vint/rules/no_asm_decl_mismatch"
	"github.com/strowk/vint/rules/no_atoi_overflow"
	"github.com/strowk/vint/rules/no_atomic_alignment_issue"
	"github.com/strowk/vint/rules/no_atomic_assign_misuse"
	"github.com/strowk/vint/rules/no_bad_lock_pattern"
	"github.com/strowk/vint/rules/no_bad_regexp_pattern"
	"github.com/strowk/vint/rules/no_bad_sort_usage"
	"github.com/strowk/vint/rules/no_bind_to_all_interfaces"
	"github.com/strowk/vint/rules/no_blank_error_assignment"
	"github.com/strowk/vint/rules/no_builtin_shadow"
	"github.com/strowk/vint/rules/no_builtin_shadow_decl"
	"github.com/strowk/vint/rules/no_capitalized_local"
	"github.com/strowk/vint/rules/no_cgi_import"
	"github.com/strowk/vint/rules/no_cgo_pointer_violation"
	"github.com/strowk/vint/rules/no_command_injection_taint"
	"github.com/strowk/vint/rules/no_commented_out_code"
	"github.com/strowk/vint/rules/no_commented_out_import"
	"github.com/strowk/vint/rules/no_conflicting_http_mux_patterns"
	"github.com/strowk/vint/rules/no_constant_parameter"
	"github.com/strowk/vint/rules/no_constant_result"
	"github.com/strowk/vint/rules/no_context_propagation_failure"
	"github.com/strowk/vint/rules/no_copied_lock"
	"github.com/strowk/vint/rules/no_deep_equal_errors"
	"github.com/strowk/vint/rules/no_defer_in_loop"
	"github.com/strowk/vint/rules/no_defer_time_misuse"
	"github.com/strowk/vint/rules/no_denied_import"
	"github.com/strowk/vint/rules/no_deprecated_hash_function"
	"github.com/strowk/vint/rules/no_direct_error_comparison"
	"github.com/strowk/vint/rules/no_doc_comment_stub"
	"github.com/strowk/vint/rules/no_duplicate_argument"
	"github.com/strowk/vint/rules/no_duplicate_branch_body"
	"github.com/strowk/vint/rules/no_duplicate_case"
	"github.com/strowk/vint/rules/no_duplicate_constants"
	"github.com/strowk/vint/rules/no_duplicate_import"
	"github.com/strowk/vint/rules/no_dynamic_errors"
	"github.com/strowk/vint/rules/no_excessive_blank_identifiers"
	"github.com/strowk/vint/rules/no_excessive_shift"
	"github.com/strowk/vint/rules/no_excessive_statements"
	"github.com/strowk/vint/rules/no_exec_command"
	"github.com/strowk/vint/rules/no_exposed_pprof"
	"github.com/strowk/vint/rules/no_file_inclusion_via_variable"
	"github.com/strowk/vint/rules/no_file_scoped_denied_import"
	"github.com/strowk/vint/rules/no_filesystem_root_serving"
	"github.com/strowk/vint/rules/no_filesystem_toctou"
	"github.com/strowk/vint/rules/no_frame_pointer_clobber"
	"github.com/strowk/vint/rules/no_hardcoded_credentials"
	"github.com/strowk/vint/rules/no_hardcoded_iv"
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
	"github.com/strowk/vint/rules/no_impossible_condition"
	"github.com/strowk/vint/rules/no_impossible_interface_assert"
	"github.com/strowk/vint/rules/no_incorrect_time_format"
	"github.com/strowk/vint/rules/no_ineffectual_assignment"
	"github.com/strowk/vint/rules/no_init_function"
	"github.com/strowk/vint/rules/no_inline_sync_once_func"
	"github.com/strowk/vint/rules/no_insecure_cookie"
	"github.com/strowk/vint/rules/no_insecure_host_key_callback"
	"github.com/strowk/vint/rules/no_insecure_random"
	"github.com/strowk/vint/rules/no_insecure_tls_config"
	"github.com/strowk/vint/rules/no_integer_overflow_conversion"
	"github.com/strowk/vint/rules/no_invalid_errors_as"
	"github.com/strowk/vint/rules/no_invalid_sort_slice_arg"
	"github.com/strowk/vint/rules/no_invalid_unsafe_pointer"
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
	"github.com/strowk/vint/rules/no_multi_line_func_break"
	"github.com/strowk/vint/rules/no_multi_line_if_break"
	"github.com/strowk/vint/rules/no_naked_return"
	"github.com/strowk/vint/rules/no_net_dial"
	"github.com/strowk/vint/rules/no_net_dial_timeout"
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
	"github.com/strowk/vint/rules/no_nil_dereference"
	"github.com/strowk/vint/rules/no_nil_func_comparison"
	"github.com/strowk/vint/rules/no_nolint_non_machine_readable"
	"github.com/strowk/vint/rules/no_nolint_parse_error"
	"github.com/strowk/vint/rules/no_nolint_without_explanation"
	"github.com/strowk/vint/rules/no_nolint_without_specific_linter"
	"github.com/strowk/vint/rules/no_non_pointer_unmarshal"
	"github.com/strowk/vint/rules/no_noop_function_call"
	"github.com/strowk/vint/rules/no_path_traversal_taint"
	"github.com/strowk/vint/rules/no_permissive_directory_permissions"
	"github.com/strowk/vint/rules/no_permissive_file_permissions"
	"github.com/strowk/vint/rules/no_permissive_os_create"
	"github.com/strowk/vint/rules/no_permissive_write_file_permissions"
	"github.com/strowk/vint/rules/no_predictable_temp_file"
	"github.com/strowk/vint/rules/no_printf_format_mismatch"
	"github.com/strowk/vint/rules/no_range_variable_alias"
	"github.com/strowk/vint/rules/no_repeated_numbers"
	"github.com/strowk/vint/rules/no_repeated_strings"
	"github.com/strowk/vint/rules/no_secret_in_serialization"
	"github.com/strowk/vint/rules/no_self_assignment"
	"github.com/strowk/vint/rules/no_serve_without_timeout"
	"github.com/strowk/vint/rules/no_short_rsa_key"
	"github.com/strowk/vint/rules/no_single_arg_append"
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
	"github.com/strowk/vint/rules/no_string_int_conversion"
	"github.com/strowk/vint/rules/no_swapped_arguments"
	"github.com/strowk/vint/rules/no_template_injection"
	"github.com/strowk/vint/rules/no_test_fatal_in_goroutine"
	"github.com/strowk/vint/rules/no_tls_conn_handshake"
	"github.com/strowk/vint/rules/no_tls_dial"
	"github.com/strowk/vint/rules/no_tls_dial_with_dialer"
	"github.com/strowk/vint/rules/no_tls_session_resumption_bypass"
	"github.com/strowk/vint/rules/no_trailing_blank_line"
	"github.com/strowk/vint/rules/no_trojan_source_bidi"
	"github.com/strowk/vint/rules/no_unallowed_import"
	"github.com/strowk/vint/rules/no_unbounded_decompression"
	"github.com/strowk/vint/rules/no_unbounded_form_parsing"
	"github.com/strowk/vint/rules/no_unbuffered_signal_channel"
	"github.com/strowk/vint/rules/no_unchecked_error"
	"github.com/strowk/vint/rules/no_unchecked_type_assertion"
	"github.com/strowk/vint/rules/no_unclosed_bodies"
	"github.com/strowk/vint/rules/no_unescaped_html_template"
	"github.com/strowk/vint/rules/no_unkeyed_literal"
	"github.com/strowk/vint/rules/no_unnecessary_conversion"
	"github.com/strowk/vint/rules/no_unnecessary_defer_lambda"
	"github.com/strowk/vint/rules/no_unnecessary_loop_var_copy"
	"github.com/strowk/vint/rules/no_unreachable_code"
	"github.com/strowk/vint/rules/no_unreachable_type_case"
	"github.com/strowk/vint/rules/no_unrecognized_directive"
	"github.com/strowk/vint/rules/no_unsafe_cors_bypass"
	"github.com/strowk/vint/rules/no_unsafe_deserialization"
	"github.com/strowk/vint/rules/no_unsafe_package"
	"github.com/strowk/vint/rules/no_unsafe_redirect_policy"
	"github.com/strowk/vint/rules/no_unused_constant"
	"github.com/strowk/vint/rules/no_unused_function"
	"github.com/strowk/vint/rules/no_unused_function_result"
	"github.com/strowk/vint/rules/no_unused_nolint"
	"github.com/strowk/vint/rules/no_unused_parameter"
	"github.com/strowk/vint/rules/no_unused_type"
	"github.com/strowk/vint/rules/no_unused_variable"
	"github.com/strowk/vint/rules/no_unused_write"
	"github.com/strowk/vint/rules/no_variable_command_execution"
	"github.com/strowk/vint/rules/no_variable_shadowing"
	"github.com/strowk/vint/rules/no_wait_group_misuse"
	"github.com/strowk/vint/rules/no_weak_crypto_hash"
	"github.com/strowk/vint/rules/no_weak_encryption_algorithm"
	"github.com/strowk/vint/rules/no_xss_taint"
	"github.com/strowk/vint/rules/no_zip_slip"
	"github.com/strowk/vint/rules/use_assignment_operator"
	"github.com/strowk/vint/rules/use_combined_append"
	"github.com/strowk/vint/rules/use_comment_spacing"
	"github.com/strowk/vint/rules/use_join_host_port"
	"github.com/strowk/vint/rules/use_matching_constant"
	"github.com/strowk/vint/rules/use_optimal_field_alignment"
	"github.com/strowk/vint/rules/use_printf_suffix"
	"github.com/strowk/vint/rules/use_simplified_bool_expr"
	"github.com/strowk/vint/rules/use_standard_codegen_comment"
	"github.com/strowk/vint/rules/use_standard_deprecation_comment"
)

var defaultRules = []lint.Rule{
	&rule.VarDeclarationsRule{},
	&rule.PackageCommentsRule{},
	&rule.DotImportsRule{},
	&rule.BlankImportsRule{},
	&rule.ExportedRule{},
	&rule.VarNamingRule{},
	&rule.IndentErrorFlowRule{},
	&rule.RangeRule{},
	&rule.ErrorfRule{},
	&rule.ErrorNamingRule{},
	&rule.ErrorStringsRule{},
	&rule.ReceiverNamingRule{},
	&rule.IncrementDecrementRule{},
	&rule.ErrorReturnRule{},
	&rule.UnexportedReturnRule{},
	&rule.TimeNamingRule{},
	&rule.ContextKeysType{},
	&rule.ContextAsArgumentRule{},
	&rule.EmptyBlockRule{},
	&rule.SuperfluousElseRule{},
	&rule.UnusedParamRule{},
	&rule.UnreachableCodeRule{},
	&rule.RedefinesBuiltinIDRule{},
	&no_unclosed_bodies.NoUnclosedBodiesRule{},
	// temporary disabled for performance research
	// &no_duplicate_code.NoDuplicateCodeRule{},
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
	&no_repeated_strings.NoRepeatedStringsRule{},
	&no_unnecessary_conversion.NoUnnecessaryConversionRule{},
	&no_unnecessary_loop_var_copy.NoUnnecessaryLoopVarCopyRule{},
	&use_printf_suffix.UsePrintfSuffixRule{},
	&no_single_arg_append.NoSingleArgAppendRule{},
	&no_asm_decl_mismatch.NoAsmDeclMismatchRule{},
	&no_self_assignment.NoSelfAssignmentRule{},
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
	&no_invalid_errors_as.NoInvalidErrorsAsRule{},
	&no_frame_pointer_clobber.NoFramePointerClobberRule{},
	&use_join_host_port.UseJoinHostPortRule{},
	&no_impossible_interface_assert.NoImpossibleInterfaceAssertRule{},
	&no_loop_closure_capture.NoLoopClosureCaptureRule{},
	&no_lost_cancel.NoLostCancelRule{},
	&no_wait_group_misuse.NoWaitGroupMisuseRule{},
	&no_nil_func_comparison.NoNilFuncComparisonRule{},
	&no_printf_format_mismatch.NoPrintfFormatMismatchRule{},
	&no_excessive_shift.NoExcessiveShiftRule{},
	&no_unbuffered_signal_channel.NoUnbufferedSignalChannelRule{},
	&no_slog_key_value_mismatch.NoSlogKeyValueMismatchRule{},
	&no_std_method_signature_mismatch.NoStdMethodSignatureMismatchRule{},
	&no_stdlib_version_mismatch.NoStdlibVersionMismatchRule{},
	&no_string_int_conversion.NoStringIntConversionRule{},
	&no_test_fatal_in_goroutine.NoTestFatalInGoroutineRule{},
	&no_malformed_test_function.NoMalformedTestFunctionRule{},
	&no_non_pointer_unmarshal.NoNonPointerUnmarshalRule{},
	&no_unreachable_code.NoUnreachableCodeRule{},
	&no_invalid_unsafe_pointer.NoInvalidUnsafePointerRule{},
	&no_unused_function_result.NoUnusedFunctionResultRule{},
	&no_unused_write.NoUnusedWriteRule{},
	&no_ineffectual_assignment.NoIneffectualAssignmentRule{},
	&no_line_too_long.NoLineTooLongRule{},
	&no_magic_number_in_argument.NoMagicNumberInArgumentRule{},
	&no_magic_number_in_assignment.NoMagicNumberInAssignmentRule{},
	&no_magic_number_in_case.NoMagicNumberInCaseRule{},
	&no_magic_number_in_condition.NoMagicNumberInConditionRule{},
	&no_magic_number_in_operation.NoMagicNumberInOperationRule{},
	&no_magic_number_in_return.NoMagicNumberInReturnRule{},

	// temporary disable while researching optimization, this one is weirdly slow
	// &no_misspelled_words.NoMisspelledWordsRule{},

	&no_naked_return.NoNakedReturnRule{},
	&no_net_dial.NoNetDialRule{},
	&no_net_dial_timeout.NoNetDialTimeoutRule{},
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
	// &no_unused_field.NoUnusedFieldRule{},
	&no_leading_blank_line.NoLeadingBlankLineRule{},
	&no_trailing_blank_line.NoTrailingBlankLineRule{},
	&no_mismatched_append_assign.NoMismatchedAppendAssignRule{},
	&use_combined_append.UseCombinedAppendRule{},
	&no_swapped_arguments.NoSwappedArgumentsRule{},
	&use_assignment_operator.UseAssignmentOperatorRule{},
	&no_noop_function_call.NoNoopFunctionCallRule{},
	&no_impossible_condition.NoImpossibleConditionRule{},
	&no_capitalized_local.NoCapitalizedLocalRule{},
	&no_unreachable_type_case.NoUnreachableTypeCaseRule{},
	&no_misplaced_default_case.NoMisplacedDefaultCaseRule{},
	&use_standard_codegen_comment.UseStandardCodegenCommentRule{},
	&use_standard_deprecation_comment.UseStandardDeprecationCommentRule{},
	&use_comment_spacing.UseCommentSpacingRule{},
	&no_duplicate_argument.NoDuplicateArgumentRule{},
	&no_duplicate_branch_body.NoDuplicateBranchBodyRule{},
	&no_duplicate_case.NoDuplicateCaseRule{},
}

var allRules = append([]lint.Rule{
	&rule.ArgumentsLimitRule{},
	&rule.CyclomaticRule{},
	&rule.FileHeaderRule{},
	&rule.ConfusingNamingRule{},
	&rule.GetReturnRule{},
	&rule.ModifiesParamRule{},
	&rule.ConfusingResultsRule{},
	&rule.DeepExitRule{},
	&rule.AddConstantRule{},
	&rule.FlagParamRule{},
	&rule.UnnecessaryStmtRule{},
	&rule.StructTagRule{},
	&rule.ModifiesValRecRule{},
	&rule.ConstantLogicalExprRule{},
	&rule.BoolLiteralRule{},
	&rule.ImportsBlocklistRule{},
	&rule.FunctionResultsLimitRule{},
	&rule.MaxPublicStructsRule{},
	&rule.RangeValInClosureRule{},
	&rule.RangeValAddress{},
	&rule.WaitGroupByValueRule{},
	&rule.AtomicRule{},
	&rule.EmptyLinesRule{},
	&rule.LineLengthLimitRule{},
	&rule.CallToGCRule{},
	&rule.DuplicatedImportsRule{},
	&rule.ImportShadowingRule{},
	&rule.BareReturnRule{},
	&rule.UnusedReceiverRule{},
	&rule.UnhandledErrorRule{},
	&rule.CognitiveComplexityRule{},
	&rule.StringOfIntRule{},
	&rule.StringFormatRule{},
	&rule.EarlyReturnRule{},
	&rule.UnconditionalRecursionRule{},
	&rule.IdenticalBranchesRule{},
	&rule.DeferRule{},
	&rule.UnexportedNamingRule{},
	&rule.FunctionLength{},
	&rule.NestedStructs{},
	&rule.UselessBreak{},
	&rule.UncheckedTypeAssertionRule{},
	&rule.TimeEqualRule{},
	&rule.TimeDateRule{},
	&rule.BannedCharsRule{},
	&rule.OptimizeOperandsOrderRule{},
	&rule.UseAnyRule{},
	&rule.DataRaceRule{},
	&rule.CommentSpacingsRule{},
	&rule.IfReturnRule{},
	&rule.RedundantImportAlias{},
	&rule.ImportAliasNamingRule{},
	&rule.EnforceMapStyleRule{},
	&rule.EnforceRepeatedArgTypeStyleRule{},
	&rule.EnforceSliceStyleRule{},
	&rule.MaxControlNestingRule{},
	&rule.CommentsDensityRule{},
	&rule.FileLengthLimitRule{},
	&rule.FilenameFormatRule{},
	&rule.RedundantBuildTagRule{},
	&rule.UseErrorsNewRule{},
	&rule.RedundantTestMainExitRule{},
	&rule.UnnecessaryFormatRule{},
	&rule.UseFmtPrintRule{},
	&rule.EnforceSwitchStyleRule{},
	&rule.IdenticalSwitchConditionsRule{},
	&rule.IdenticalIfElseIfConditionsRule{},
	&rule.IdenticalIfElseIfBranchesRule{},
	&rule.IdenticalSwitchBranchesRule{},
	&rule.UselessFallthroughRule{},
	&rule.PackageDirectoryMismatchRule{},
	&rule.UseWaitGroupGoRule{},
	&rule.UnsecureURLSchemeRule{},
	&rule.InefficientMapLookupRule{},
	&rule.ForbiddenCallInWgGoRule{},
	&rule.UnnecessaryIfRule{},
	&rule.EpochNamingRule{},
	&rule.UseSlicesSort{},
	&rule.PackageNamingRule{},
	&no_atomic_alignment_issue.NoAtomicAlignmentIssueRule{},
	&no_blank_error_assignment.NoBlankErrorAssignmentRule{},
	&no_specific_function_call.NoSpecificFunctionCallRule{},
	&no_deep_equal_errors.NoDeepEqualErrorsRule{},
	&no_duplicate_constants.NoDuplicateConstantsRule{},
	&no_unchecked_type_assertion.NoUncheckedTypeAssertionRule{},
	&no_repeated_numbers.NoRepeatedNumbersRule{},
	&use_matching_constant.UseMatchingConstantRule{},
	&use_optimal_field_alignment.UseOptimalFieldAlignmentRule{},
	&no_conflicting_http_mux_patterns.NoConflictingHttpMuxPatternsRule{},
	&no_invalid_sort_slice_arg.NoInvalidSortSliceArgRule{},
	&no_nil_dereference.NoNilDereferenceRule{},
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
	&no_commented_out_code.NoCommentedOutCodeRule{},
	&no_commented_out_import.NoCommentedOutImportRule{},
	&no_defer_in_loop.NoDeferInLoopRule{},
	&no_unnecessary_defer_lambda.NoUnnecessaryDeferLambdaRule{},
	&no_doc_comment_stub.NoDocCommentStubRule{},
	&no_duplicate_import.NoDuplicateImportRule{},
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
	switch name {
	case "imports-blacklist":
		return "imports-blocklist"
	default:
		return name
	}
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
