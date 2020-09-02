package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func add(x int, y int) int {
	return x + y
}

func server() {
	PORT := ":8080"
	l, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer l.Close()

	c, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		netData, err := bufio.NewReader(c).ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		if strings.TrimSpace(string(netData)) == "STOP" {
			fmt.Println("Exiting TCP server!")
			return
		}

		fmt.Print("-> ", string(netData))
		t := time.Now()
		myTime := t.Format(time.RFC3339) + "\n"
		c.Write([]byte(myTime))
	}
}

func listen(c net.Conn) {

	for {
		reader := bufio.NewReader(os.Stdin)

		fmt.Print(">> ")
		text, _ := reader.ReadString('\n')
		fmt.Fprintf(c, text+"\n")
		message, _ := bufio.NewReader(c).ReadString('\n')
		fmt.Print("->: " + message)
		if strings.TrimSpace(string(text)) == "STOP" {
			fmt.Println("TCP client exiting...")
			return
		}
	}
}

func main() {
	go server()

	fmt.Println("listener starts")
	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("enter IP to connect to >> ")
		text, _ := reader.ReadString('\n')
		CONNECT := strings.TrimSuffix(text, "\n") + ":8080"
		fmt.Println(CONNECT)
		c, err := net.Dial("tcp", CONNECT)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Starting Client")
			listen(c)
			break
		}
	}
}
