package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

var in *os.File = os.Stdin

// Cli starts the program for the given node and outputs data to the given
// io.writer
func Cli(output io.Writer, node kademlia.Node) {
	fmt.Fprintln(out, "Starting CLI...")
	reader := bufio.NewReader(in)

	for {
		input, _ := reader.ReadString('\n')
		trimInput := strings.TrimSpace(input)

		if trimInput == "\n" || trimInput == "" {
			continue

		} else {
			commands := strings.Fields(trimInput)

			Commands(output, &node, commands)

		}

	}
}
