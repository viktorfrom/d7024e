package kademlia

import (
	"errors"
	"net"
	"time"

	log "github.com/sirupsen/logrus"
)

type Message struct {
	reciever Contact
	rpc      RPC
	err      error
}

type Client struct {
	ip   string
	send chan Message
	resp chan Message
}

func InitClient() Client {
	client := Client{}
	client.ip = client.GetLocalIP()
	client.send = make(chan Message)
	client.resp = make(chan Message)
	go client.sendRPC()
	return client
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

func (client *Client) sendRPC() {
	for {
		message := <-client.send
		rpc := message.rpc
		reciever := message.reciever

		if rpc.TargetID == nil || rpc.SenderID == nil {
			log.Warn(errNoID)
			client.resp <- Message{Contact{}, RPC{}, errors.New(errNoID)}
			continue
		}

		sendRPCID := *rpc.ID
		readBuffer := make([]byte, 1024)

		msg, err := MarshalRPC(rpc)
		if err != nil {
			log.Warn(err)
			client.resp <- Message{Contact{}, RPC{}, err}
			continue
		}

		sendAddr, err := net.ResolveUDPAddr(udpNetwork, reciever.Address)
		if err != nil {
			log.Warn(err)
			client.resp <- Message{Contact{}, RPC{}, err}
			continue
		}

		conn, err := net.DialUDP(udpNetwork, nil, sendAddr)
		if err != nil {
			log.Warn(err)
			client.resp <- Message{Contact{}, RPC{}, err}
			continue
		}
		defer conn.Close()

		conn.SetDeadline(time.Now().Add(timeout))
		conn.SetReadDeadline(time.Now().Add(timeout))

		_, err = conn.Write(msg)
		if err != nil {
			log.Warn(err)
			client.resp <- Message{Contact{}, RPC{}, err}
			continue
		}

		bytesRead, receiveAddr, err := conn.ReadFromUDP(readBuffer)
		if err != nil {
			log.Warn(err)
			client.resp <- Message{Contact{}, RPC{}, err}
			continue
		}

		if bytesRead == 0 {
			log.Warn(errNoReply)
			client.resp <- Message{Contact{}, RPC{}, errors.New(errNoReply)}
			continue
		}

		if receiveAddr.String() != sendAddr.String() {
			log.Warn(errDiffAddr)
			client.resp <- Message{Contact{}, RPC{}, errors.New(errDiffAddr)}
			continue
		}

		reply, err := UnmarshalRPC(readBuffer[0:bytesRead])
		if err != nil {
			log.Warn(err)
			client.resp <- Message{Contact{}, RPC{}, err}
			continue
		}

		if sendRPCID != *reply.ID {
			log.Warn(errDiffID)
			client.resp <- Message{Contact{}, RPC{}, errors.New(errDiffID)}
			continue
		}

		client.resp <- Message{reciever, *reply, nil}
	}
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

	client.send <- Message{*contact, *rpc, nil}

	resp := <-client.resp

	if resp.err != nil {
		log.Warn(resp.err)
		return nil, resp.err
	}

	return &resp.rpc, resp.err
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

	client.send <- Message{*contact, *rpc, nil}
	resp := <-client.resp

	if resp.err != nil {
		log.Warn(resp.err)
		return nil, resp.err
	}

	return &resp.rpc, resp.err
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

	client.send <- Message{*contact, *rpc, nil}
	resp := <-client.resp

	if resp.err != nil {
		log.Warn(resp.err)
		return nil, resp.err
	}

	return &resp.rpc, resp.err
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
	// rpc, err := client.sendRPC(contact, Store, sender.ID, contact.ID, payload)
	rpc, _ := NewRPC(Store, sender.ID.String(), contact.ID.String(), payload)

	client.send <- Message{*contact, *rpc, nil}
	resp := <-client.resp

	if resp.err != nil {
		log.Warn(resp.err)
		return nil, resp.err
	}

	return &resp.rpc, resp.err
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
