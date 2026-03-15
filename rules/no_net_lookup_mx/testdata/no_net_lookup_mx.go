package fixtures

import (
	"context"
	"fmt"
	"net"
)

func badNetLookupMX() {
	mxs, err := net.LookupMX("example.com") // MATCH /net.LookupMX does not accept a context; use (*net.Resolver).LookupMX instead/
	if err != nil {
		panic(err)
	}
	for _, mx := range mxs {
		fmt.Println(mx.Host, mx.Pref)
	}
}

func goodNetLookupMX() {
	ctx := context.Background()
	r := net.Resolver{}
	mxs, err := r.LookupMX(ctx, "example.com")
	if err != nil {
		panic(err)
	}
	for _, mx := range mxs {
		fmt.Println(mx.Host, mx.Pref)
	}
}
