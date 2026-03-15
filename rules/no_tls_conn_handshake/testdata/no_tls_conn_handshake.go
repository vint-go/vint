package fixtures

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
)

func badTlsConnHandshake() {
	conn, err := net.Dial("tcp", "example.com:443")
	if err != nil {
		panic(err)
	}
	tlsConn := tls.Client(conn, &tls.Config{})
	err = tlsConn.Handshake() // MATCH /(*tls.Conn).Handshake does not accept a context; use (*tls.Conn).HandshakeContext instead/
	if err != nil {
		panic(err)
	}
	defer tlsConn.Close()
	fmt.Println("Handshake completed")
}

func goodTlsConnHandshakeContext() {
	ctx := context.Background()
	conn, err := net.Dial("tcp", "example.com:443")
	if err != nil {
		panic(err)
	}
	tlsConn := tls.Client(conn, &tls.Config{})
	err = tlsConn.HandshakeContext(ctx)
	if err != nil {
		panic(err)
	}
	defer tlsConn.Close()
	fmt.Println("Handshake completed")
}
