package main

import (
	"bufio"
	"os"
	"strings"
)

func Cli() {
	reader := bufio.NewReader(os.Stdin)

	for {
		input, _ := reader.ReadString('\n')
		trimInput := strings.TrimSpace(input)

		if trimInput == "\n" || trimInput == "" {
			continue
		} else {
			Commands(trimInput)
		}

	}
}
