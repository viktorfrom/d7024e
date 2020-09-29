package kademlia

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewNetwork(t *testing.T) {
	node := Node{}
	network := NewNetwork(&node)

	assert.NotNil(t, network)
	assert.Equal(t, node, *network.kademlia)
}

func TestGetLocalIp(t *testing.T) {
	node := Node{}
	network := NewNetwork(&node)
	assert.NotNil(t, network.ip)
	assert.NotEqual(t, "", network.GetLocalIP())
	assert.Equal(t, network.ip, network.GetLocalIP())
}

func TestUpdateRoutingTable(t *testing.T) {
	pingMsg := pingMsg
	payload := Payload{&pingMsg, nil, nil}

	c := NewContact(NewNodeID("1111111400000000000000000000000000000000"), "localhost:8002")

	target := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")
	rpc, _ := NewRPC(Ping, target.ID.String(), payload)

	node := Node{}
	node.RT = NewRoutingTable(c)
	network := NewNetwork(&node)

	assert.Equal(t, []Contact(nil), node.RT.FindClosestContacts(c.ID, 5))
	network.updateRoutingTable(rpc, "10.0.8.1")
	target.CalcDistance(c.ID)
	assert.Equal(t, []Contact{target}, node.RT.FindClosestContacts(c.ID, 5))
}

func TestHandleIncomingPing(t *testing.T) {
	node := Node{}
	network := NewNetwork(&node)
	pingMsg := pingMsg
	payload := Payload{&pingMsg, nil, nil}
	target := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")

	rpc, _ := NewRPC(Ping, target.ID.String(), payload)
	r, err := network.handleIncomingPingRPC(rpc)
	assert.Equal(t, r, rpc)
	assert.Nil(t, err)

	r, err = network.handleIncomingPingRPC(nil)
	assert.Nil(t, r)
	assert.Equal(t, errors.New(errNilRPC), err)
}

func TestHandleIncomingFindNode(t *testing.T) {
	node := Node{}
	network := NewNetwork(&node)
	payload := Payload{nil, nil, []Contact{}}
	target := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")

	rpc, err := NewRPC(FindNode, target.ID.String(), payload)

	r, err := network.handleIncomingFindNodeRPC(rpc)
	assert.Equal(t, r, rpc)
	assert.Equal(t, errors.New(errNoContact), err)

	_, err = network.handleIncomingFindNodeRPC(nil)
	assert.Equal(t, errors.New(errNilRPC), err)
}

func TestHandleIncomingRPCS(t *testing.T) {
	node := Node{}
	c := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")
	node.RT = NewRoutingTable(c)
	network := NewNetwork(&node)
	pingMsg := pingMsg
	payload := Payload{nil, &pingMsg, nil}

	orgRPC, _ := NewRPC(Ping, "1111111100000000000000000000000000000000", payload)
	rpc, err := network.handleIncomingRPCS(orgRPC, "10.0.8.3:8080")

	assert.Equal(t, orgRPC, rpc)
	assert.Equal(t, OK, *rpc.Type)
	assert.Nil(t, err)

	storeRPC, _ := NewRPC(Store, "1111111100000000000000000000000000000000", Payload{nil, nil, []Contact{}})
	_, err = network.handleIncomingRPCS(storeRPC, "10.0.8.3:8080")
	assert.Error(t, err)

	nodeRPC, _ := NewRPC(FindNode, "1111111100000000000000000000000000000000", Payload{nil, nil, []Contact{}})
	_, err = network.handleIncomingRPCS(nodeRPC, "10.0.8.3:8080")
	assert.Error(t, err)

	valueRPC, _ := NewRPC(FindValue, "1111111100000000000000000000000000000000", Payload{nil, nil, []Contact{}})
	_, err = network.handleIncomingRPCS(valueRPC, "10.0.8.3:8080")
	assert.Error(t, err)

	wrongRPC, _ := NewRPC(OK, "1111111100000000000000000000000000000000", Payload{nil, nil, []Contact{}})
	_, err = network.handleIncomingRPCS(wrongRPC, "10.0.8.3:8080")
	assert.Error(t, err)
}

func TestPingError(t *testing.T) {
	node := Node{}
	c := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")
	node.RT = NewRoutingTable(c)
	network := NewNetwork(&node)

	_, err := network.SendPingMessage(nil, &c)
	assert.Error(t, err)
}

func TestFindNodeError(t *testing.T) {
	network := Network{}
	_, err := network.SendFindContactMessage(nil, nil)
	assert.Error(t, err)
}

func TestFindValueError(t *testing.T) {
	node := Node{}
	c := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")
	node.RT = NewRoutingTable(c)
	network := NewNetwork(&node)

	_, err := network.SendFindDataMessage(&c, nil, "testkey")
	assert.Error(t, err)
}

func TestStoreError(t *testing.T) {
	network := Network{}
	_, err := network.SendStoreMessage(nil, nil, "key", "value")
	assert.Error(t, err)
}

func TestHandleIncomingStoreError(t *testing.T) {
	storeType := Store

	network := Network{}
	_, err := network.handleIncomingStoreRPC(nil)
	assert.Error(t, err)

	rpc := RPC{&storeType, nil, nil, nil}
	_, err = network.handleIncomingStoreRPC(&rpc)
	assert.Error(t, err)

	payload := Payload{nil, nil, []Contact{}}
	rpc = RPC{&storeType, &payload, nil, nil}
	_, err = network.handleIncomingStoreRPC(&rpc)
	assert.Error(t, err)
}

func TestHandleIncomingFindValue(t *testing.T) {
	findValue := FindValue
	network := Network{}

	_, err := network.handleIncomingFindValueRPC(nil)
	assert.Equal(t, errors.New(errNilRPC), err)

	payload := Payload{nil, nil, []Contact{}}
	rpc := RPC{&findValue, &payload, nil, nil}
	_, err = network.handleIncomingFindValueRPC(&rpc)
	assert.Equal(t, errors.New(errBadKeyValue), err)
}

func TestListenErrors(t *testing.T) {
	network := Network{}
	// Port 1 is reserved and can never be used
	err := network.Listen("127.0.0.1", "1")
	assert.Error(t, err)
}

func TestSendRPCNoNetwork(t *testing.T) {
	timeout = 0 * time.Second

	network := Network{}
	nodeID := NewNodeID("00000000000000000000000000000000FFFFFFFF")
	c := NewContact(nodeID, "10.0.8.1:8080")

	payload := Payload{nil, nil, []Contact{}}
	_, err := network.sendRPC(&c, Ping, nodeID, payload)
	assert.Error(t, err)
}
