package kademlia

import (
	"errors"
	"log"
	"net"

	. "github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

const (
	udpNetwork  string = "udp4"
	pingMsg     string = "PING"
	pongMsg     string = "PONG"
	errNoReply  string = "did not receive a reply"
	errDiffAddr string = "receive address not same as send address"
)

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
type Network struct {
}

func Listen(ip string, port int) {
	listenAddr := &net.UDPAddr{IP: net.ParseIP(ip), Port: port, Zone: ""}
	readBuffer := make([]byte, 1024)

	conn, err := net.ListenUDP(udpNetwork, listenAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	for {
		bytesRead, receiveAddr, err := conn.ReadFromUDP(readBuffer)
		if err != nil {
			log.Fatal(err)
		}

		receivedMsg := string(readBuffer[0:bytesRead])
		if receivedMsg == pingMsg {
			conn.WriteToUDP([]byte(pongMsg), receiveAddr)
		}
	}
}

// SendPingMessage pings a contact and returns the response. Returns an
// error if the contact fails to respond.
func (network *Network) SendPingMessage(contact *Contact) (*string, error) {
	readBuffer := make([]byte, 1024)
	sendAddr, err := net.ResolveUDPAddr(udpNetwork, contact.Address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP(udpNetwork, nil, sendAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(pingMsg))
	if err != nil {
		return nil, err
	}

	bytesRead, receiveAddr, err := conn.ReadFromUDP(readBuffer)
	if err != nil {
		return nil, err
	}
	response := string(readBuffer[0:bytesRead])

	if bytesRead == 0 {
		return nil, errors.New(errNoReply)
	}

	if receiveAddr.String() != sendAddr.String() {
		return nil, errors.New(errDiffAddr)
	}

	return &response, nil
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
