package fixtures

import (
	"context"
	"net"
	"time"
)

func badNetDialTimeout() {
	conn, err := net.DialTimeout("tcp", "example.com:80", 5*time.Second) // MATCH /net.DialTimeout does not accept a context; use (*net.Dialer).DialContext with Timeout instead/
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}

func goodNetDialerDialContext() {
	ctx := context.Background()
	d := net.Dialer{
		Timeout: 5 * time.Second,
	}
	conn, err := d.DialContext(ctx, "tcp", "example.com:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}
