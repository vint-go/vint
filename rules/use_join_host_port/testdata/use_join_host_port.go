package fixtures

import (
	"fmt"
	"net"
)

func connectBadIntPort(host string, port int) (net.Conn, error) {
	// Bad: does not work with IPv6 addresses
	addr := fmt.Sprintf("%s:%d", host, port) // MATCH /use net.JoinHostPort instead of fmt.Sprintf for host:port construction to support IPv6/
	return net.Dial("tcp", addr)
}

func connectBadStringPort(host string, port string) (net.Conn, error) {
	// Bad: does not work with IPv6 addresses
	addr := fmt.Sprintf("%s:%s", host, port) // MATCH /use net.JoinHostPort instead of fmt.Sprintf for host:port construction to support IPv6/
	return net.Dial("tcp", addr)
}

func connectGood(host string, port string) (net.Conn, error) {
	// Good: net.JoinHostPort handles IPv6 correctly
	addr := net.JoinHostPort(host, port)
	return net.Dial("tcp", addr)
}

func formatOther() string {
	// Not a host:port pattern - should not trigger
	return fmt.Sprintf("%s:%s:%s", "a", "b", "c")
}

func formatDifferent() string {
	// Different format string - should not trigger
	return fmt.Sprintf("%s=%d", "key", 42)
}

func formatNonLiteral(format string) string {
	// Non-literal format string - should not trigger
	return fmt.Sprintf(format, "host", 80)
}
