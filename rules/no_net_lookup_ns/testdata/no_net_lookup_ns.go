package fixtures

import (
	"context"
	"fmt"
	"net"
)

func badNetLookupNS() {
	nss, err := net.LookupNS("example.com") // MATCH /net.LookupNS does not accept a context; use (*net.Resolver).LookupNS instead/
	if err != nil {
		panic(err)
	}
	for _, ns := range nss {
		fmt.Println(ns.Host)
	}
}

func goodNetLookupNS() {
	ctx := context.Background()
	r := net.Resolver{}
	nss, err := r.LookupNS(ctx, "example.com")
	if err != nil {
		panic(err)
	}
	for _, ns := range nss {
		fmt.Println(ns.Host)
	}
}
