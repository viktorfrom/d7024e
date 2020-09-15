package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func Commands(Input string) {
	switch Input {
	case "put":
		fmt.Println("Put")
	case "p":
		fmt.Println("Put")
	case "get":
		fmt.Println("Get")
	case "g":
		fmt.Println("Get")
	case "exit":
		Exit()
	case "e":
		Exit()
	case "--help":
		Help()
	case "--h":
		Help()
	case "--version":
		Help()
	case "--v":
		Help()
	default:
		fmt.Println("Invalid command!")
	}
}

func Put() {
	fmt.Println("Put")

}

func Get() {
	fmt.Println("Get")

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
