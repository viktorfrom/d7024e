package main

import (
	"bufio"
	"os"
	"strings"

	"github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

func Cli(node kademlia.Node) {
	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')
		trimInput := strings.TrimSpace(input)

		if trimInput == "\n" || trimInput == "" {
			continue
		} else {
			commands := strings.Fields(trimInput)

			Commands(node, commands)

		}

	}
}
