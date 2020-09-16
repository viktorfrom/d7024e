package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

func Commands(node kademlia.Kademlia, commands []string) {

	switch commands[0] {
	case "put":
		if len(commands) == 2 {
			Put(node, commands[1])
		}
	case "p":
		if len(commands) == 2 {
			Put(node, commands[1])
		}
	case "ping":
		if len(commands) == 2 {
			Ping(node, commands[1])
		}
	case "get":
		if len(commands) == 2 {
			Get(node, commands[1])
		}
	case "g":
		if len(commands) == 2 {
			Get(node, commands[1])
		}
	case "exit":
		Exit()
	case "e":
		Exit()
	case "help":
		Help()
	case "h":
		Help()
	case "version":
		Help()
	case "v":
		Help()
	default:
		fmt.Println("Invalid command!")
	}
}

func Put(node kademlia.Kademlia, input string) {
	// Convert string to []byte
	data := []byte(input)
	node.Store(data)
}

func Ping(node kademlia.Kademlia, input string) {
	node.Ping()
}

func Get(node kademlia.Kademlia, hash string) {
	node.LookupData(hash)
}

func Exit() {
	os.Exit(3)
}

func Help() {
	content, err := ioutil.ReadFile("prompt.txt")
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	text := string(content)
	fmt.Println(text)

}
