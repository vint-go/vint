package fixtures

import (
	"context"
	"crypto/tls"
)

func badTlsDial() {
	conn, err := tls.Dial("tcp", "example.com:443", &tls.Config{}) // MATCH /tls.Dial does not accept a context; use (*tls.Dialer).DialContext instead/
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}

func goodTlsDialContext() {
	ctx := context.Background()
	d := tls.Dialer{
		Config: &tls.Config{},
	}
	conn, err := d.DialContext(ctx, "tcp", "example.com:443")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}
