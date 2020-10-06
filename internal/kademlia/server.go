package kademlia

import (
	"errors"
	"net"
	"strconv"
	"strings"
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

const (
	errNoReply        string = "did not receive a reply"
	errDiffAddr       string = "receive address not same as send address"
	errDiffID         string = "RPC ID was different"
	errNilRPC         string = "RPC struct is nil"
	errInvalidRPCType string = "RPC type is invalid"
	errNoContact      string = "no contact was given"
	errNoTargetID     string = "no TargetID given"
	errNoID           string = "no ID given"
	errBadKeyValue    string = "bad or no key or value given"
	errNoRPCPayload   string = "no RPC payload given"
)

// the time before a RPC call times out
var timeout = 10 * time.Second

type Server struct {
	kademlia *Node
	ip       string
}

// NewNetwork initializes the network and sets the local IP address
func InitServer(kademlia *Node) Server {
	server := Server{}
	server.kademlia = kademlia
	server.ip = server.GetLocalIP()
	return server
}

// GetLocalIP returns the IP of the Node in the Docker Network
func (server *Server) GetLocalIP() string {
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
func (server *Server) Listen(port string) error {
	portAsInt, _ := strconv.Atoi(port)
	listenAddr := &net.UDPAddr{IP: net.ParseIP(server.ip), Port: portAsInt, Zone: ""}

	conn, err := net.ListenUDP(udpNetwork, listenAddr)
	if err != nil {
		log.Error(err)
		return err
	}
	defer conn.Close()

	err = server.handleUDP(conn)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}

func (server *Server) handleUDP(conn *net.UDPConn) error {
	readBuffer := make([]byte, UDPReadBufferSize)

	for {
		bytesRead, receiveAddr, err := conn.ReadFromUDP(readBuffer)
		if err != nil {
			log.Warn(err)
			continue
		} else if bytesRead == 0 {
			continue
		}

		rpc, err := UnmarshalRPC(readBuffer[0:bytesRead])
		if err != nil {
			log.Warn(err)
			continue
		}

		rpc, err = server.handleIncomingRPCS(rpc, strings.Split(receiveAddr.String(), ":")[0])
		if err != nil {
			log.Warn(err)
			continue
		}
		data, err := MarshalRPC(*rpc)

		conn.WriteToUDP(data, receiveAddr)
	}
}

func (server *Server) handleIncomingRPCS(rpc *RPC, receiveAddr string) (*RPC, error) {
	var err error
	var retRPC *RPC
	switch *rpc.Type {
	case Ping:
		retRPC, err = server.handleIncomingPingRPC(rpc)
	case Store:
		retRPC, err = server.handleIncomingStoreRPC(rpc)
	case FindNode:
		retRPC, err = server.handleIncomingFindNodeRPC(rpc)
	case FindValue:
		retRPC, err = server.handleIncomingFindValueRPC(rpc)
	default:
		return rpc, errors.New(errInvalidRPCType)
	}

	if err != nil {
		log.Warn(err)
		return nil, err
	}

	server.updateRoutingTable(rpc, receiveAddr)
	*rpc.Type = OK
	*rpc.SenderID = server.kademlia.RT.GetMeID().String()

	return retRPC, nil
}

func (server *Server) updateRoutingTable(rpc *RPC, senderIP string) {
	sender := NewNodeID(*rpc.SenderID)
	contact := NewContact(sender, senderIP+DefaultPort)
	server.kademlia.RT.AddContact(contact)
}

func (server *Server) handleIncomingPingRPC(rpc *RPC) (*RPC, error) {
	if rpc == nil {
		return nil, errors.New(errNilRPC)
	}

	return rpc, nil
}

func (server *Server) handleIncomingStoreRPC(rpc *RPC) (*RPC, error) {
	err := checkNilRPCPayload(rpc)
	if err != nil {
		return nil, err
	}

	key := rpc.Payload.Key
	value := rpc.Payload.Value
	if key == nil || value == nil {
		return nil, errors.New(errBadKeyValue)
	}

	server.kademlia.insertLocalStore(*key, *value)

	return rpc, nil
}

func (server *Server) handleIncomingFindNodeRPC(rpc *RPC) (*RPC, error) {
	err := checkNilRPCPayload(rpc)
	if err != nil {
		return nil, err
	}

	if rpc.TargetID == nil {
		return nil, errors.New(errNoTargetID)
	}

	targetID := NewNodeID(*rpc.TargetID)
	contacts := server.kademlia.RT.FindClosestContacts(targetID, BucketSize)

	payload := Payload{nil, nil, contacts}
	rpc.Payload = &payload

	return rpc, nil
}

func (server *Server) handleIncomingFindValueRPC(rpc *RPC) (*RPC, error) {
	err := checkNilRPCPayload(rpc)
	if err != nil {
		log.Warn(err)
		return nil, err
	}

	if rpc.TargetID == nil {
		log.Warn(errNoTargetID)
		return nil, errors.New(errNoTargetID)
	}

	key := rpc.Payload.Key
	if key == nil {
		log.Warn(errBadKeyValue)
		return nil, errors.New(errBadKeyValue)
	}

	value := server.kademlia.searchLocalStore(*key)
	// If no value is found - return k closest
	if value == nil {
		return server.handleIncomingFindNodeRPC(rpc)
	}

	rpc.Payload.Value = value
	return rpc, nil
}
