package fixtures

import (
	"context"
	"net"
)

func badNetLookupTXT() {
	txts, err := net.LookupTXT("example.com") // MATCH /net.LookupTXT does not accept a context; use (*net.Resolver).LookupTXT instead/
	if err != nil {
		panic(err)
	}
	_ = txts
}

func goodNetLookupTXT() {
	ctx := context.Background()
	r := net.Resolver{}
	txts, err := r.LookupTXT(ctx, "example.com")
	if err != nil {
		panic(err)
	}
	_ = txts
}
