package fixtures

import (
	"context"
	"crypto/tls"
	"net"
	"time"
)

func badTlsDialWithDialer() {
	dialer := &net.Dialer{
		Timeout: 5 * time.Second,
	}
	// tls.DialWithDialer does not accept a context
	conn, err := tls.DialWithDialer(dialer, "tcp", "example.com:443", &tls.Config{}) // MATCH /tls.DialWithDialer does not accept a context; use (*tls.Dialer).DialContext with NetDialer instead/
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}

func goodTlsDialerDialContext() {
	ctx := context.Background()
	d := tls.Dialer{
		NetDialer: &net.Dialer{
			Timeout: 5 * time.Second,
		},
		Config: &tls.Config{},
	}
	conn, err := d.DialContext(ctx, "tcp", "example.com:443")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}
