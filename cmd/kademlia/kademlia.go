package main

import (
	"fmt"

	"github.com/viktorfrom/d7024e-kademlia/internal/network"
)

func main() {
	fmt.Println("Hello, Arch!")

	// setup UDP server on port 8080
	go network.Listen("8080")

	n := network.Network{}
	n.SendPingMessage(nil)
}
