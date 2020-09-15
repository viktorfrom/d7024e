package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

func Commands(node kademlia.Kademlia, Input string) {
	switch Input {
	case "put":
		Put(node, Input)
	case "p":
		Put(node, Input)
	case "get":
		Get(node, Input)
	case "g":
		Get(node, Input)
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

func Put(node kademlia.Kademlia, Input string) {
	// Convert string to []byte
	data := []byte(Input)
	node.Store(data)
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
