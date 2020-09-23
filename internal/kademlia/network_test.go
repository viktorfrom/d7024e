package kademlia

import (
	"errors"
	"testing"

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
}
