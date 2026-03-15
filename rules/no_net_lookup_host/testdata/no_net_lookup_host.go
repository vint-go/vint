package fixtures

import (
	"context"
	"net"
)

func badNetLookupHost() {
	addrs, err := net.LookupHost("example.com") // MATCH /net.LookupHost does not accept a context; use (*net.Resolver).LookupHost instead/
	if err != nil {
		panic(err)
	}
	_ = addrs
}

func goodNetLookupHost() {
	ctx := context.Background()
	r := net.Resolver{}
	addrs, err := r.LookupHost(ctx, "example.com")
	if err != nil {
		panic(err)
	}
	_ = addrs
}
