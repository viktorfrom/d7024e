package kademlia

import (
	"errors"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

// the time before a RPC call times out
var timeout = 10 * time.Second

//Message represents the data sent through the channels in the Client
type Message struct {
	receiver Contact
	rpc      RPC
	err      error
}

type Client struct {
	ip   string
	send chan Message // channel for sending messages from the client to a server
	resp chan Message // channel for the responses from the server to the client
}

//InitClient sets up and returns a client object
func InitClient() Client {
	client := Client{}
	client.ip = client.GetLocalIP()
	client.send = make(chan Message)
	client.resp = make(chan Message)

	return client
}

// Start starts the client which runs in a goroutine
func (client *Client) Start() {
	go func() {
		for {
			err := client.sendRPC()
			if err != nil {
				client.resp <- Message{Contact{}, RPC{}, err}
				log.Warn(err)
			}
		}
	}()
}

// GetLocalIP returns the IP of the Node in the Docker Network
func (client *Client) GetLocalIP() string {
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

func (client *Client) sendRPC() error {
	message := <-client.send
	rpc := message.rpc
	receiver := message.receiver

	if rpc.TargetID == nil || rpc.SenderID == nil {
		return errors.New(errNoID)
	}

	readBuffer := make([]byte, 1024)

	msg, err := MarshalRPC(rpc)
	if err != nil {
		return err
	}

	sendAddr, conn, err := client.createConnection(receiver.Address)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(msg)
	if err != nil {
		return err
	}

	reply, err := client.getRPCReply(conn, readBuffer, sendAddr)
	if err != nil {
		return err
	}

	// validate the reply
	if *rpc.ID != *reply.ID {
		return errors.New(errDiffID)
	}

	client.resp <- Message{receiver, *reply, nil}
	return nil
}

// createConnection returns an UDPAddr and UDPConn to the given address
func (client *Client) createConnection(addr string) (*net.UDPAddr, *net.UDPConn, error) {
	sendAddr, err := net.ResolveUDPAddr(udpNetwork, addr)
	if err != nil {
		return nil, nil, err
	}

	conn, err := net.DialUDP(udpNetwork, nil, sendAddr)
	if err != nil {
		return sendAddr, nil, err
	}

	conn.SetDeadline(time.Now().Add(timeout))
	conn.SetReadDeadline(time.Now().Add(timeout))

	return sendAddr, conn, nil

}

func (client *Client) getRPCReply(conn *net.UDPConn, readBuffer []byte, sendAddr *net.UDPAddr) (*RPC, error) {
	bytesRead, receiveAddr, err := conn.ReadFromUDP(readBuffer)
	if err != nil {
		return nil, err
	}

	reply, err := client.getRPCFromBuffer(bytesRead, receiveAddr.String(), sendAddr.String(), readBuffer)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (client *Client) getRPCFromBuffer(bytesRead int, receiveAddr, sendAddr string, readBuffer []byte) (*RPC, error) {
	if bytesRead == 0 {
		return nil, errors.New(errNoReply)
	}

	if receiveAddr != sendAddr {
		return nil, errors.New(errDiffAddr)
	}

	reply, err := UnmarshalRPC(readBuffer[0:bytesRead])
	if err != nil {
		return nil, err
	}

	return reply, nil
}

func (client *Client) sendMessage(rpc *RPC, contact *Contact) (*RPC, error) {
	client.send <- Message{*contact, *rpc, nil}
	resp := <-client.resp

	if resp.err != nil {
		log.Warn(resp.err)
		return nil, resp.err
	}

	return &resp.rpc, resp.err
}

// SendPingMessage sends a PING RPC to the `contact` and returns an acknowledgement. `sender` is needed in
// case the receiving node needs information about the node who sent the RPC. Returns an error
// if the contact fails to respond or any argument is invalid.
func (client *Client) SendPingMessage(contact *Contact, sender *Contact) (*RPC, error) {
	err := checkNilContacts(contact, sender)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	pingMsg := pingMsg
	payload := Payload{nil, &pingMsg, nil}
	rpc, _ := NewRPC(Ping, sender.ID.String(), contact.ID.String(), payload)

	return client.sendMessage(rpc, contact)
}

// SendFindContactMessage sends a FIND_NODE RPC to `contact`. `sender` is needed in case the receiving
// node needs information about the node who sent the RPC. `targetID` is the NodeID which is targeted in this RPC.
// Returns an error if the contact fails to respond or any argument is invalid.
func (client *Client) SendFindContactMessage(contact, sender *Contact, targetID *NodeID) (*RPC, error) {
	err := checkNilContacts(contact, sender)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	payload := Payload{nil, nil, []Contact{}}
	rpc, _ := NewRPC(FindNode, sender.ID.String(), targetID.String(), payload)

	return client.sendMessage(rpc, contact)
}

// SendFindDataMessage sends a FIND_VALUE RPC to `contact` looking for the value belonging to `key`. If the
// value is found it will return the stored value otherwise the contacts `k` closest nodes will return.
// Note that `key` is the hash of the value, it is used as a TargetID internally because they share the same
// ID space. Returns an error if the contact fails to respond or any argument is invalid.
func (client *Client) SendFindDataMessage(contact, sender *Contact, key string) (*RPC, error) {
	err := checkNilContacts(contact, sender)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	targetID := NewNodeID(key)
	payload := Payload{&key, nil, nil}
	rpc, _ := NewRPC(FindValue, sender.ID.String(), targetID.String(), payload)

	return client.sendMessage(rpc, contact)
}

// SendStoreMessage sends a STORE RPC to `contact` with a given `key`, `value`. `sender` is the node that sends this
// RPC. Note that `key` is the hash of `value`. Returns an error if the contact fails to respond or any argument is invalid.
func (client *Client) SendStoreMessage(contact *Contact, sender *Contact, key string, value string) (*RPC, error) {
	err := checkNilContacts(contact, sender)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	payload := Payload{&key, &value, nil}
	rpc, _ := NewRPC(Store, sender.ID.String(), contact.ID.String(), payload)

	return client.sendMessage(rpc, contact)
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
