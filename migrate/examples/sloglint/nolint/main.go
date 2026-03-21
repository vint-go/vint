package example

import "log/slog"

func logUser() {
	slog.Info("a user has logged in", "user_id", 42, slog.String("ip_address", "192.0.2.0")) //nolint:sloglint
}
