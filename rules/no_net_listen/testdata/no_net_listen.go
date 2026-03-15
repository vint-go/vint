package fixtures

import (
	"context"
	"net"
)

func badNetListen() {
	ln, err := net.Listen("tcp", ":8080") // MATCH /net.Listen does not accept a context; use (*net.ListenConfig).Listen instead/
	if err != nil {
		panic(err)
	}
	defer ln.Close()
}

func goodNetListen() {
	ctx := context.Background()
	lc := net.ListenConfig{}
	ln, err := lc.Listen(ctx, "tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer ln.Close()
}
