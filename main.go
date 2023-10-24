package main

import (
	"flag"
)

var listenAddress = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")

func main() {
	s := NewServer()
	s.Run()
}
