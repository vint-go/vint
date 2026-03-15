package fixtures

import (
	"context"
	"net"
)

func badNetLookupCNAME() {
	cname, err := net.LookupCNAME("www.example.com") // MATCH /net.LookupCNAME does not accept a context; use (*net.Resolver).LookupCNAME instead/
	if err != nil {
		panic(err)
	}
	_ = cname
}

func goodNetLookupCNAME() {
	ctx := context.Background()
	r := net.Resolver{}
	cname, err := r.LookupCNAME(ctx, "www.example.com")
	if err != nil {
		panic(err)
	}
	_ = cname
}
