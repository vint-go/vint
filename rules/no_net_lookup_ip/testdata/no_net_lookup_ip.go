package fixtures

import (
	"context"
	"net"
)

func badNetLookupIp() {
	ips, err := net.LookupIP("example.com") // MATCH /net.LookupIP does not accept a context; use (*net.Resolver).LookupIPAddr instead/
	if err != nil {
		panic(err)
	}
	_ = ips
}

func goodNetLookupIp() {
	ctx := context.Background()
	r := net.Resolver{}
	ips, err := r.LookupIPAddr(ctx, "example.com")
	if err != nil {
		panic(err)
	}
	_ = ips
}
