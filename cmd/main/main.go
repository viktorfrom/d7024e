package main

import (
	"fmt"
	"io"
	"os"

	"github.com/viktorfrom/d7024e-kademlia/cmd/api"
	"github.com/viktorfrom/d7024e-kademlia/cmd/cli"
	"github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

var out io.Writer = os.Stdout

func main() {
	fmt.Fprintln(out, "Booting Kademlia....")

	node := kademlia.Node{}
	node.InitNode()

	server := kademlia.InitServer(&node)
	go server.Listen("8080")
	go api.API(out, node)

	cli.Cli(out, node)
}
