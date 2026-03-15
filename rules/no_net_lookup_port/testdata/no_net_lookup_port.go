package fixtures

import (
	"context"
	"net"
)

func badNetLookupPort() {
	port, err := net.LookupPort("tcp", "http") // MATCH /net.LookupPort does not accept a context; use (*net.Resolver).LookupPort instead/
	if err != nil {
		panic(err)
	}
	_ = port
}

func goodNetLookupPort() {
	ctx := context.Background()
	r := net.Resolver{}
	port, err := r.LookupPort(ctx, "tcp", "http")
	if err != nil {
		panic(err)
	}
	_ = port
}
