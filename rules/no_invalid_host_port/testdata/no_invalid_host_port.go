package fixtures

import (
	"net"
	"net/http"
)

func badMissingColon() {
	http.ListenAndServe("localhost8080", nil) // MATCH /invalid host:port pair "localhost8080" passed to http.ListenAndServe/
}

func badNetListen() {
	net.Listen("tcp", "localhost8080") // MATCH /invalid host:port pair "localhost8080" passed to net.Listen/
}

func badNetDial() {
	net.Dial("tcp", "localhost8080") // MATCH /invalid host:port pair "localhost8080" passed to net.Dial/
}

func badNoPort() {
	http.ListenAndServe("localhost", nil) // MATCH /invalid host:port pair "localhost" passed to http.ListenAndServe/
}

// Valid examples below: these should not trigger failures

func goodCorrectFormat() {
	http.ListenAndServe(":8080", nil)
}

func goodHostAndPort() {
	http.ListenAndServe("localhost:8080", nil)
}

func goodEmptyString() {
	http.ListenAndServe("", nil)
}

func goodNetListen() {
	net.Listen("tcp", ":8080")
}

func goodNetDial() {
	net.Dial("tcp", "localhost:8080")
}

func goodVariable(addr string) {
	http.ListenAndServe(addr, nil)
}

func goodIPv6() {
	net.Dial("tcp", "[::1]:8080")
}

func goodColonPort() {
	http.ListenAndServe(":0", nil)
}
