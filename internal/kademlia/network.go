package kademlia

import (
	"errors"
	"net"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	// DefaultPort Default port to listen on
	DefaultPort string = ":8080"
	// UDPReadBufferSize Size of the UDP read buffer
	UDPReadBufferSize int = 1024
)

const (
	udpNetwork   string = "udp4"
	pingMsg      string = "PING"
	errNoReply   string = "did not receive a reply"
	errDiffAddr  string = "receive address not same as send address"
	errDiffID    string = "rpc ID was different"
	errNilRPC    string = "rpc struct is nil"
	errNoContact string = "no contact was given"
)

// the time before a RPC call times out
const timeout = 10 * time.Second

type Network struct {
	kademlia *Node
	ip       string
}

// NewNetwork initializes the network and sets the local IP address
func NewNetwork(kademlia *Node) Network {
	network := Network{}
	network.kademlia = kademlia
	network.ip = network.GetLocalIP()
	return network
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

// Listen listens on the given ip and port. If it fails to connect
// an error will be returned.
func (network *Network) Listen(ip string, port string) error {
	portAsInt, _ := strconv.Atoi(port)
	listenAddr := &net.UDPAddr{IP: net.ParseIP(ip), Port: portAsInt, Zone: ""}

	conn, err := net.ListenUDP(udpNetwork, listenAddr)
	if err != nil {
		log.Error(err)
		return err
	}
	defer conn.Close()

	err = network.handleIncomingRPCS(conn)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (network *Network) handleIncomingRPCS(conn *net.UDPConn) error {
	readBuffer := make([]byte, UDPReadBufferSize)

	for {
		bytesRead, receiveAddr, err := conn.ReadFromUDP(readBuffer)
		if err != nil {
			log.Warn(err)
			continue
		}

		rpc, err := UnmarshalRPC(readBuffer[0:bytesRead])
		if err != nil {
			log.Warn(err)
			continue
		}

		switch *rpc.Type {
		case Ping:
			rpc, err = network.handleIncomingPingRPC(rpc)
		case Store:
			rpc, err = network.handleIncomingStoreRPC(rpc)
		case FindNode:
			rpc, err = network.handleIncomingFindNodeRPC(rpc)
		case FindValue:
			rpc, err = network.handleIncomingFindValueRPC(rpc)
		default:
			continue
		}

		if err != nil {
			log.Warn(err)
			continue
		}

		network.updateRoutingTable(rpc, receiveAddr.String())
		*rpc.Type = OK
		data, _ := MarshalRPC(*rpc)
		conn.WriteToUDP(data, receiveAddr)
	}
}

func (network *Network) updateRoutingTable(rpc *RPC, senderIP string) {
	sender := NewNodeID(*rpc.SenderID)
	contact := NewContact(sender, senderIP+DefaultPort)
	network.kademlia.RT.AddContact(contact)
}

func (network *Network) handleIncomingPingRPC(rpc *RPC) (*RPC, error) {
	if rpc == nil {
		return nil, errors.New(errNilRPC)
	}

	return rpc, nil
}

func (network *Network) handleIncomingStoreRPC(rpc *RPC) (*RPC, error) {
	return rpc, nil
}

func (network *Network) handleIncomingFindNodeRPC(rpc *RPC) (*RPC, error) {
	if len(rpc.Payload.Contacts) == 0 {
		return rpc, errors.New(errNoContact)
	}

	contacts := network.kademlia.RT.FindClosestContacts(rpc.Payload.Contacts[0].ID, BucketSize)

	payload := Payload{nil, contacts}
	*rpc.Payload = payload
	*rpc.Type = OK

	return rpc, nil
}

func (network *Network) handleIncomingFindValueRPC(rpc *RPC) (*RPC, error) {
	return rpc, nil
}

func (network *Network) sendRPC(contact *Contact, rpcType RPCType, senderID *NodeID, payload Payload) (*RPC, error) {
	rpc, _ := NewRPC(rpcType, senderID.String(), payload)
	sendRPCID := *rpc.ID
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

	if sendRPCID != *reply.ID {
		return nil, errors.New(errDiffID)
	}

	return reply, nil
}

// SendPingMessage pings a contact and returns the response. `sender` is needed in case the receiving node needs information
// about the node who sent the RPC. Returns an error if the contact fails to respond.
func (network *Network) SendPingMessage(contact *Contact, sender *Contact) (*RPC, error) {
	pingMsg := pingMsg
	payload := Payload{&pingMsg, nil}
	rpc, err := network.sendRPC(contact, Ping, sender.ID, payload)

	return rpc, err
}

// SendFindContactMessage sends a FindNode RPC to contact. `sender` is needed in case the receiving node needs information
// about the node who sent the RPC. Returns an error if the contact fails to respond.
func (network *Network) SendFindContactMessage(contact *Contact, sender *Contact) (*RPC, error) {
	payload := Payload{nil, []Contact{*contact}}
	rpc, err := network.sendRPC(contact, FindNode, sender.ID, payload)

	return rpc, err
}

// SendFindDataMessage TODO
func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

// SendStoreMessage TODO
func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
