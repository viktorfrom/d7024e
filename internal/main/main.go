package main

import (
	"fmt"
	"io"
	"os"

	cli "github.com/viktorfrom/d7024e-kademlia/cmd"
	"github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

var out io.Writer = os.Stdout

func main() {
	fmt.Fprintln(out, "Booting Kademlia....")

	node := kademlia.Node{}
	node.InitNode()

	cli.Cli(out, node)
}
