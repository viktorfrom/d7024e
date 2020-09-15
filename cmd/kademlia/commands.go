package main

import (
	"fmt"
	"io/ioutil"
	"log"
)

func Commands(Input string) {
	switch Input {
	case "put":
		fmt.Println("Good morning!")
	case "get":
		fmt.Println("Good afternoon!")
	case "exit":
		fmt.Println("Good afternoon!")
	case "help":
		Help()
	case "version":
		fmt.Println("Good afternoon!")
	default:
		fmt.Println("Good evening!")
	}
}

func Put(Input string) {
	fmt.Println("Put")

}

func Get(Input string) {
	fmt.Println("Get")

}

func Exit(Input string) {
	fmt.Println("Exit")

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
