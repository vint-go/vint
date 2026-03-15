package fixtures

import (
	"context"
	"net"
)

func badNetDial() {
	conn, err := net.Dial("tcp", "example.com:80") // MATCH /net.Dial does not accept a context; use (*net.Dialer).DialContext instead/
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}

func goodNetDialContext() {
	ctx := context.Background()
	d := net.Dialer{}
	conn, err := d.DialContext(ctx, "tcp", "example.com:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}
