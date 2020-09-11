package main

import (
	"fmt"

	"github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

func main() {
	fmt.Println("Hello, Arch!")

	node := kademlia.Kademlia{}
	node.InitNode(kademlia.NewRandomKademliaID())

	fmt.Println(node.RT.ME.String())
}
