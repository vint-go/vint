package fixtures

import (
	"fmt"
	"net"
)

func connectURLIntPort(host string, port int) string {
	// Bad: URL with scheme and host:port - does not work with IPv6 addresses
	return fmt.Sprintf("http://%s:%d/path", host, port) // MATCH /use net.JoinHostPort instead of fmt.Sprintf for host:port construction to support IPv6/
}

func connectURLStringPort(host string, port string) string {
	// Bad: URL with scheme and host:port - does not work with IPv6 addresses
	return fmt.Sprintf("https://%s:%s/path", host, port) // MATCH /use net.JoinHostPort instead of fmt.Sprintf for host:port construction to support IPv6/
}

func connectURLMinimal(host string, port int) string {
	// Bad: URL with scheme and host:port, no trailing path
	return fmt.Sprintf("http://%s:%d", host, port) // MATCH /use net.JoinHostPort instead of fmt.Sprintf for host:port construction to support IPv6/
}

func connectURLCustomScheme(host string, port string) string {
	// Bad: custom scheme URL with host:port
	return fmt.Sprintf("myapp+tcp://%s:%s", host, port) // MATCH /use net.JoinHostPort instead of fmt.Sprintf for host:port construction to support IPv6/
}

func connectBareIntPort(host string, port int) (net.Conn, error) {
	// Bad: bare host:port does not work with IPv6 addresses
	addr := fmt.Sprintf("%s:%d", host, port) // MATCH /use net.JoinHostPort instead of fmt.Sprintf for host:port construction to support IPv6/
	return net.Dial("tcp", addr)
}

func connectBareStringPort(host string, port string) (net.Conn, error) {
	// Bad: bare host:port does not work with IPv6 addresses
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
