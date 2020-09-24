package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

func Commands(node kademlia.Node, commands []string) {

	switch commands[0] {
	case "put":
		if len(commands) == 2 {
			Put(node, commands[1])
		} else {
			fmt.Println("No argument!")
		}
	case "p":
		if len(commands) == 2 {
			Put(node, commands[1])
		} else {
			fmt.Println("No argument!")
		}
	case "ping":
		if len(commands) == 2 {
			Ping(node, commands[1])
		}
	case "get":
		if len(commands) == 2 {
			Get(node, commands[1])
		} else {
			fmt.Println("No argument!")
		}
	case "g":
		if len(commands) == 2 {
			Get(node, commands[1])
		} else {
			fmt.Println("No argument!")
		}
	case "t":
		// c := node.RT.GetMe()
		c := kademlia.NewContact(kademlia.NewNodeID("00000000000000100b5e0038281912513b2f5751"), "10.0.8.9")
		c.CalcDistance(node.RT.GetMeID())
		node.NodeLookup(&c)
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

func Put(node kademlia.Node, input string) {
	node.StoreValue(input)
}

func Ping(node kademlia.Node, input string) {
	node.Ping()
}

func Get(node kademlia.Node, hash string) {
	node.FindValue(hash)
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
