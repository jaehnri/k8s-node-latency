package main

import "k8s-node-latency/internal/client"

func main() {
	c := client.NewClient()
	c.Run()
}
