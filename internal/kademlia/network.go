package kademlia

import (
	"errors"
	"net"
	"strconv"
	"time"

	. "github.com/viktorfrom/d7024e-kademlia/internal/rpc"
)

const (
	udpNetwork  string = "udp4"
	pingMsg     string = "PING"
	pongMsg     string = "PONG"
	errNoReply  string = "did not receive a reply"
	errDiffAddr string = "receive address not same as send address"
	errDiffID   string = "rpc ID was different"
)

// the time before a RPC call times out
const timeout = 10 * time.Second

type Network struct{}

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

// Listen listens on the given ip and port. If it fails to connect
// an error will be returned.
func (network *Network) Listen(ip string, port string) error {
	portAsInt, _ := strconv.Atoi(port)
	listenAddr := &net.UDPAddr{IP: net.ParseIP(ip), Port: portAsInt, Zone: ""}

	conn, err := net.ListenUDP(udpNetwork, listenAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = handleIncomingRPCS(conn)
	if err != nil {
		return err
	}

	return nil
}

func handleIncomingRPCS(conn *net.UDPConn) error {
	readBuffer := make([]byte, 1024)

	for {
		bytesRead, receiveAddr, err := conn.ReadFromUDP(readBuffer)
		if err != nil {
			return err
		}

		rpc, err := UnmarshalRPC(readBuffer[0:bytesRead])
		if err != nil {
			return err
		}

		*rpc.Type = OK
		data, _ := MarshalRPC(*rpc)
		conn.WriteToUDP(data, receiveAddr)
	}
}

func (network *Network) sendRPC(contact *Contact, rpcType RPCType, data []byte) (*RPC, error) {
	rpc, _ := NewRPC(rpcType, data)
	sendID := *rpc.ID
	readBuffer := make([]byte, 1024)

	msg, err := MarshalRPC(*rpc)
	if err != nil {
		return nil, err
	}

	sendAddr, err := net.ResolveUDPAddr(udpNetwork, contact.Address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP(udpNetwork, nil, sendAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(timeout))
	conn.SetReadDeadline(time.Now().Add(timeout))

	_, err = conn.Write(msg)
	if err != nil {
		return nil, err
	}

	bytesRead, receiveAddr, err := conn.ReadFromUDP(readBuffer)
	if err != nil {
		return nil, err
	}

	if bytesRead == 0 {
		return nil, errors.New(errNoReply)
	}

	if receiveAddr.String() != sendAddr.String() {
		return nil, errors.New(errDiffAddr)
	}

	reply, err := UnmarshalRPC(readBuffer[0:bytesRead])
	if err != nil {
		return nil, err
	}

	if sendID != *reply.ID {
		return nil, errors.New(errDiffID)
	}

	return reply, nil
}

// SendPingMessage pings a contact and returns the response. Returns an
// error if the contact fails to respond.
func (network *Network) SendPingMessage(contact *Contact) (*RPC, error) {
	rpc, err := network.sendRPC(contact, Ping, []byte(pingMsg))
	if err != nil {
		return nil, err
	}

	return rpc, nil
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
