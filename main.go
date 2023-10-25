package main

import (
	"flag"
)

var (
	tcpListenAddress  = flag.String("tcp-listen-address", ":3000", "The address to listen on for TCP requests.")
	httpListenAddress = flag.String("http-listen-address", ":8080", "The address to listen on for HTTP requests.")
)

func main() {
	s := NewServer(*tcpListenAddress, *httpListenAddress)
	s.Run()
}
