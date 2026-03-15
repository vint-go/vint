package fixtures

import (
	"context"
	"net"
)

func badNetLookupSRV() {
	_, srvs, err := net.LookupSRV("xmpp-server", "tcp", "example.com") // MATCH /net.LookupSRV does not accept a context; use (*net.Resolver).LookupSRV instead/
	if err != nil {
		panic(err)
	}
	for _, srv := range srvs {
		_ = srv.Target
	}
}

func goodNetLookupSRV() {
	ctx := context.Background()
	r := net.Resolver{}
	_, srvs, err := r.LookupSRV(ctx, "xmpp-server", "tcp", "example.com")
	if err != nil {
		panic(err)
	}
	for _, srv := range srvs {
		_ = srv.Target
	}
}
