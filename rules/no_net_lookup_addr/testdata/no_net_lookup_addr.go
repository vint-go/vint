package fixtures

import (
	"context"
	"net"
)

func badNetLookupAddr() {
	names, err := net.LookupAddr("127.0.0.1") // MATCH /net.LookupAddr does not accept a context; use (*net.Resolver).LookupAddr instead/
	if err != nil {
		panic(err)
	}
	_ = names
}

func goodNetLookupAddr() {
	ctx := context.Background()
	r := net.Resolver{}
	names, err := r.LookupAddr(ctx, "127.0.0.1")
	if err != nil {
		panic(err)
	}
	_ = names
}
