package fixtures

import (
	"context"
	"net"
)

func badNetListenPacket() {
	conn, err := net.ListenPacket("udp", ":8080") // MATCH /net.ListenPacket does not accept a context; use (*net.ListenConfig).ListenPacket instead/
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}

func goodNetListenPacket() {
	ctx := context.Background()
	lc := net.ListenConfig{}
	conn, err := lc.ListenPacket(ctx, "udp", ":8080")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
}
