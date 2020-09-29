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
	udpNetwork string = "udp4"
	pingMsg    string = "PING"
)

// Network error messages
const (
	errNoReply        string = "did not receive a reply"
	errDiffAddr       string = "receive address not same as send address"
	errDiffID         string = "RPC ID was different"
	errNilRPC         string = "RPC struct is nil"
	errInvalidRPCType string = "RPC type is invalid"
	errNoContact      string = "no contact was given"
	errBadKeyValue    string = "bad or no key or value given"
	errNoRPCPayload   string = "no RPC payload given"
)

// the time before a RPC call times out
var timeout = 10 * time.Second

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

// Listen listens on the given ip and port for incoming RPC requests.
//If it fails to connect an error will be returned.
func (network *Network) Listen(ip string, port string) error {
	portAsInt, _ := strconv.Atoi(port)
	listenAddr := &net.UDPAddr{IP: net.ParseIP(ip), Port: portAsInt, Zone: ""}

	conn, err := net.ListenUDP(udpNetwork, listenAddr)
	if err != nil {
		log.Error(err)
		return err
	}
	defer conn.Close()

	err = network.handleUDP(conn)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (network *Network) handleUDP(conn *net.UDPConn) error {
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

		rpc, err = network.handleIncomingRPCS(rpc, receiveAddr.String())
		data, err := MarshalRPC(*rpc)

		conn.WriteToUDP(data, receiveAddr)
	}
}

func (network *Network) handleIncomingRPCS(rpc *RPC, receiveAddr string) (*RPC, error) {
	var err error
	var retRPC *RPC
	switch *rpc.Type {
	case Ping:
		retRPC, err = network.handleIncomingPingRPC(rpc)
	case Store:
		retRPC, err = network.handleIncomingStoreRPC(rpc)
	case FindNode:
		retRPC, err = network.handleIncomingFindNodeRPC(rpc)
	case FindValue:
		retRPC, err = network.handleIncomingFindValueRPC(rpc)
	default:
		return nil, errors.New(errInvalidRPCType)
	}

	if err != nil {
		log.Warn(err)
		return nil, err
	}

	network.updateRoutingTable(rpc, receiveAddr)
	*rpc.Type = OK
	*rpc.SenderID = network.kademlia.RT.GetMeID().String()

	return retRPC, nil
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
	err := checkNilRPCPayload(rpc)
	if err != nil {
		return nil, err
	}

	key := rpc.Payload.Key
	value := rpc.Payload.Value
	if key == nil || value == nil {
		return nil, errors.New(errBadKeyValue)
	}

	network.kademlia.insertLocalStore(*key, *value)

	return rpc, nil
}

func (network *Network) handleIncomingFindNodeRPC(rpc *RPC) (*RPC, error) {
	err := checkNilRPCPayload(rpc)
	if err != nil {
		return nil, err
	}

	if len(rpc.Payload.Contacts) == 0 {
		return rpc, errors.New(errNoContact)
	}

	contacts := network.kademlia.RT.FindClosestContacts(rpc.Payload.Contacts[0].ID, BucketSize)

	payload := Payload{nil, nil, contacts}
	*rpc.Payload = payload

	return rpc, nil
}

func (network *Network) handleIncomingFindValueRPC(rpc *RPC) (*RPC, error) {
	err := checkNilRPCPayload(rpc)
	if err != nil {
		return nil, err
	}

	key := rpc.Payload.Key
	if key == nil {
		return nil, errors.New(errBadKeyValue)
	}

	value := network.kademlia.searchLocalStore(*key)

	// If no value is found - return k closest
	if value == nil {
		return network.handleIncomingFindNodeRPC(rpc)
	}

	*rpc.Payload.Value = *value
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

// SendPingMessage sends a PING RPC to a contact and returns the response. `sender` is needed in
// case the receiving node needs information about the node who sent the RPC. Returns an error
// if the contact fails to respond.
func (network *Network) SendPingMessage(contact *Contact, sender *Contact) (*RPC, error) {
	err := checkNilContacts(contact, sender)
	if err != nil {
		return nil, err
	}

	pingMsg := pingMsg
	payload := Payload{nil, &pingMsg, nil}
	rpc, err := network.sendRPC(contact, Ping, sender.ID, payload)

	return rpc, err
}

// SendFindContactMessage sends a FIND_NODE RPC to contact. `sender` is needed in case the receiving
// node needs information about the node who sent the RPC. Returns an error if the contact fails to respond.
func (network *Network) SendFindContactMessage(contact *Contact, sender *Contact) (*RPC, error) {
	err := checkNilContacts(contact, sender)
	if err != nil {
		return nil, err
	}

	payload := Payload{nil, nil, []Contact{*contact}}
	rpc, err := network.sendRPC(contact, FindNode, sender.ID, payload)

	return rpc, err
}

// SendFindDataMessage sends a FIND_VALUE RPC to contact looking for the value belonging to `key`. If the
// value is found it will return the stored value otherwise the contacts `k` closest nodes will return.
// Returns an error if the contact fails to respond.
func (network *Network) SendFindDataMessage(contact *Contact, sender *Contact, key string) (*RPC, error) {
	err := checkNilContacts(contact, sender)
	if err != nil {
		return nil, err
	}

	payload := Payload{&key, nil, nil}
	rpc, err := network.sendRPC(contact, FindValue, sender.ID, payload)

	return rpc, err
}

// SendStoreMessage sends a STORE RPC to contact with a given `key`, `value`. Returns an error if the contact
// fails to respond.
func (network *Network) SendStoreMessage(contact *Contact, sender *Contact, key string, value string) (*RPC, error) {
	err := checkNilContacts(contact, sender)
	if err != nil {
		return nil, err
	}

	payload := Payload{&key, &value, nil}
	rpc, err := network.sendRPC(contact, Store, sender.ID, payload)

	return rpc, err
}

func checkNilRPCPayload(rpc *RPC) error {
	if rpc == nil {
		return errors.New(errNilRPC)
	}

	if rpc.Payload == nil {
		return errors.New(errNoRPCPayload)
	}
	return nil
}

func checkNilContacts(contact *Contact, sender *Contact) error {
	if contact == nil && sender == nil {
		return errors.New(errNoContact + ": contact & sender")
	} else if contact == nil {
		return errors.New(errNoContact + ": contact")
	} else if sender == nil {
		return errors.New(errNoContact + ": sender")
	}
	return nil
}
