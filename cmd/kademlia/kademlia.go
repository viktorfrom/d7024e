package main

import (
	"fmt"

	"github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

func main() {
	fmt.Println("Booting Kademlia....")

	node := kademlia.Kademlia{}
	go node.InitNode(kademlia.NewRandomKademliaID())

	Cli(node)
}
