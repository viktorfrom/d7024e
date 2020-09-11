package kademlia

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Network TODO
type Network struct {
}

// InitNetwork TODO
func (network *Network) InitNetwork(ip string, port string) {
	network.Listen(ip, port)
}

// GetLocalIP returns the IP of the Node in the Docker Network
func (network *Network) GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}

// Listen Start UDP server
func (network *Network) Listen(ip string, port string) {
	fmt.Println("Starting server")
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

	for {
		n, addr, err := connection.ReadFromUDP(buffer)
		fmt.Print("-> ", string(buffer[0:n]))

		var data []byte

		if strings.TrimSpace(string(buffer[0:n])) == "PING" {
			data = []byte(string("PONG"))
		} else if strings.TrimSpace(string(buffer[0:n])) == "FIND_NODE" {
			data = []byte(string(ip + ";" + port + ";SuperRandomID"))
		}

		fmt.Printf("data: %s\n", string(data))
		_, err = connection.WriteToUDP(data, addr)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}

// setup for the network client
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

// sendMessage sends a stream of data to a UDP server
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
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("enter IP to connect to >> ")
		ip, _ := reader.ReadString('\n')
		CONNECT := strings.TrimSuffix(ip, "\n") + ":8080"
		fmt.Print("Command >> ")
		command, _ := reader.ReadString('\n')

		data := []byte(strings.TrimSuffix(command, "\n"))

		network.sendMessage(CONNECT, data)
	}
}

// SendFindContactMessage TODO
func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

// SendFindDataMessage TODO
func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

// SendStoreMessage TODO
func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
