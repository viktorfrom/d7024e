package main

import (
	"fmt"

	"github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

func main() {
	fmt.Println("Booting Kademlia....")

	node := kademlia.Node{}
	node.InitNode(kademlia.NewRandomNodeID())

	Cli(node)
}
