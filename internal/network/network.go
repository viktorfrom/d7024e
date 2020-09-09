package network

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strings"
	"time"

	. "github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

// Network ....
type Network struct {
}

// Listen Start UDP server
func Listen(ip string, port string) {

	PORT := ":" + port
	s, err := net.ResolveUDPAddr("udp4", PORT)
	if err != nil {
		fmt.Println(err)
		return
	}

	connection, err := net.ListenUDP("udp4", s)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()
	buffer := make([]byte, 1024)
	rand.Seed(time.Now().Unix())

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		fmt.Print("-> ", string(buffer[0:n]))

		var data []byte

		if strings.TrimSpace(string(buffer[0:n])) == "PING" {
			data = []byte(string("PONG"))
		} else if strings.TrimSpace(string(buffer[0:n])) == "FIND_NODE" {
			data = []byte(string("IP;8080;SuperRandomID"))
		}

		fmt.Printf("data: %s\n", string(data))
		_, err = connection.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

func (network *Network) initNetworkClient(ip string) net.UDPConn {
	fmt.Println(ip)
	s, err := net.ResolveUDPAddr("udp4", ip)
	if err != nil {
		fmt.Println(err)
	}
	c, err := net.DialUDP("udp4", nil, s)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("The UDP server is %s\n", c.RemoteAddr().String())
	}
	return *c
}

func (network *Network) sendMessage(ip string, data []byte) {
	c := network.initNetworkClient(ip)
	_, err := c.Write(data)

	if err != nil {
		fmt.Println(err)
	}

	if strings.TrimSpace(string(data)) == "STOP" {
		fmt.Println("Exiting UDP client!")
		return
	}

	buffer := make([]byte, 1024)
	n, _, err := c.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Reply: %s\n", string(buffer[0:n]))
	c.Close()
}

// SendPingMessage Temporary UDP client sending example packages to an inputed IP address
func (network *Network) SendPingMessage(contact *Contact) {
	// TODO

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("enter IP to connect to >> ")
		ip, _ := reader.ReadString('\n')
		CONNECT := strings.TrimSuffix(ip, "\n") + ":8080"
		fmt.Print(">> ")
		command, _ := reader.ReadString('\n')

		data := []byte(strings.TrimSuffix(command, "\n"))

		network.sendMessage(CONNECT, data)
	}
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
