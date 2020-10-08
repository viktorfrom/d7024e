package kademlia

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNetwork(t *testing.T) {
	node := Node{}
	network := InitServer(&node)

	assert.NotNil(t, network)
	assert.Equal(t, node, *network.kademlia)
}

func TestGetLocalIp(t *testing.T) {
	node := Node{}
	network := InitServer(&node)
	assert.NotNil(t, network.ip)
	assert.NotEqual(t, "", network.GetLocalIP())
	assert.Equal(t, network.ip, network.GetLocalIP())
}

func TestUpdateRoutingTable(t *testing.T) {
	pingMsg := pingMsg
	payload := Payload{&pingMsg, nil, nil}

	c := NewContact(NewNodeID("1111111400000000000000000000000000000000"), "localhost:8002")

	target := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")
	rpc, _ := NewRPC(Ping, target.ID.String(), "", payload)

	node := Node{}
	node.RT = NewRoutingTable(c)
	network := InitServer(&node)

	assert.Equal(t, []Contact(nil), node.RT.FindClosestContacts(c.ID, 5))
	network.updateRoutingTable(rpc, "10.0.8.1")
	target.CalcDistance(c.ID)
	assert.Equal(t, []Contact{target}, node.RT.FindClosestContacts(c.ID, 5))
}

func TestHandleIncomingPing(t *testing.T) {
	node := Node{}
	network := InitServer(&node)
	pingMsg := pingMsg
	payload := Payload{&pingMsg, nil, nil}
	target := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")

	rpc, _ := NewRPC(Ping, target.ID.String(), "", payload)
	r, err := network.handleIncomingPingRPC(rpc)
	assert.Equal(t, r, rpc)
	assert.Nil(t, err)

	r, err = network.handleIncomingPingRPC(nil)
	assert.Nil(t, r)
	assert.Equal(t, errors.New(errNilRPC), err)

}

func TestHandleIncomingRPCS(t *testing.T) {
	node := Node{}
	c := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")
	node.RT = NewRoutingTable(c)
	network := InitServer(&node)
	pingMsg := pingMsg
	payload := Payload{nil, &pingMsg, nil}

	orgRPC, _ := NewRPC(Ping, "1111111100000000000000000000000000000000", "00000000000000000000000000000000FFFFFFFF", payload)
	rpc, err := network.handleIncomingRPCS(orgRPC, "10.0.8.3:8080")

	assert.Equal(t, orgRPC, rpc)
	assert.Equal(t, OK, *rpc.Type)
	assert.Nil(t, err)

	storeRPC, _ := NewRPC(Store, "1111111100000000000000000000000000000000", "00000000000000000000000000000000FFFFFFFF", Payload{nil, nil, []Contact{}})
	_, err = network.handleIncomingRPCS(storeRPC, "10.0.8.3:8080")
	assert.Error(t, err)

	valueRPC, _ := NewRPC(FindValue, "1111111100000000000000000000000000000000", "00000000000000000000000000000000FFFFFFFF", Payload{nil, nil, []Contact{}})
	_, err = network.handleIncomingRPCS(valueRPC, "10.0.8.3:8080")
	assert.Error(t, err)

	wrongRPC, _ := NewRPC(OK, "1111111100000000000000000000000000000000", "00000000000000000000000000000000FFFFFFFF", Payload{nil, nil, []Contact{}})
	_, err = network.handleIncomingRPCS(wrongRPC, "10.0.8.3:8080")
	assert.Error(t, err)
}

func TestHandleIncomingRPCsFindValue(t *testing.T) {
	node := Node{}
	c := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")
	node.RT = NewRoutingTable(c)
	network := InitServer(&node)

	nodeRPC, _ := NewRPC(FindNode, "1111111100000000000000000000000000000000", "00000100000000000000000000000000FFFFFFFF", Payload{nil, nil, []Contact{}})
	_, err := network.handleIncomingRPCS(nodeRPC, "10.0.8.3:8080")
	assert.Nil(t, err)
}

func TestListenReservedPortError(t *testing.T) {
	network := Server{}
	// Port 1 is reserved and can never be used so should always throw error
	err := network.Listen("1")
	assert.Error(t, err)
}

// The below tests, tests pure functionality
// don't move anything above to below here before refactoring the tests!
func TestIncomingFindNodeFindsCorrectContact(t *testing.T) {
	senderID := "00000000000000000000000000000000FFFFFFFF"
	targetID := "1111111100000000000000000000000000000000"

	senderAddress := "10.0.8.1:8080"
	targetAddress := "10.0.8.2:8080"

	node := Node{}
	c := NewContact(NewNodeID(senderID), senderAddress)
	node.RT = NewRoutingTable(c)
	network := InitServer(&node)
	payload := Payload{nil, nil, []Contact{}}

	target := NewContact(NewNodeID(targetID), targetAddress)
	node.RT.AddContact(target)
	rpc, _ := NewRPC(FindNode, senderID, targetID, payload)
	rpc, err := network.handleIncomingFindNodeRPC(rpc)
	assert.Nil(t, err)
	assert.Equal(t, rpc.Payload.Contacts[0].ID, target.ID)
}

func TestIncomingFindNodeBadInput(t *testing.T) {
	node := Node{}
	c := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")
	node.RT = NewRoutingTable(c)
	network := InitServer(&node)
	payload := Payload{nil, nil, []Contact{}}

	rpc, err := NewRPC(FindNode, "00000000000000000000000000000000FFFFFFFF", "1111111100000000000000000000000000000000", payload)

	_, err = network.handleIncomingFindNodeRPC(rpc)
	assert.Nil(t, err)

	_, err = network.handleIncomingFindNodeRPC(nil)
	assert.Equal(t, errors.New(errNilRPC), err)

	rpc.TargetID = nil
	_, err = network.handleIncomingFindNodeRPC(rpc)
	assert.Equal(t, errors.New(errNoTargetID), err)
}

func TestIncomingFindValueNoKeyInPayload(t *testing.T) {
	targetID := "1111111100000000000000000000000000000000"

	findValue := FindValue
	network := Server{}

	_, err := network.handleIncomingFindValueRPC(nil)
	assert.Equal(t, errors.New(errNilRPC), err)

	payload := Payload{nil, nil, []Contact{}}
	rpc := RPC{&findValue, &payload, nil, nil, &targetID}
	_, err = network.handleIncomingFindValueRPC(&rpc)
	assert.Equal(t, errors.New(errBadKeyValue), err)
}

func TestIncomingFindValueReturnClosestContacts(t *testing.T) {
	node := Node{}
	c := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")
	node.RT = NewRoutingTable(c)
	network := InitServer(&node)

	key := "1111111100000000000000000000000000000000"
	payload := Payload{&key, nil, []Contact{}}
	rpc, _ := NewRPC(FindValue, "00000000000000000000000000000000FFFFFFFF", "1111111100000000000000000000000000000000", payload)
	rpc, _ = network.handleIncomingFindValueRPC(rpc)

	assert.Nil(t, rpc.Payload.Value)
	assert.Equal(t, []Contact(nil), rpc.Payload.Contacts)
}

func TestIncomingFindValueFoundValue(t *testing.T) {
	node := Node{nil, Client{}, make(map[string]string)}
	c := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")
	node.RT = NewRoutingTable(c)
	network := InitServer(&node)

	key := "1111111100000000000000000000000000000000"
	value := "hello"
	payload := Payload{&key, nil, []Contact{}}
	node.insertLocalStore(key, value)
	rpc, _ := NewRPC(FindValue, "00000000000000000000000000000000FFFFFFFF", "1111111100000000000000000000000000000000", payload)
	rpc, err := network.handleIncomingFindValueRPC(rpc)
	assert.Nil(t, err)
	assert.Equal(t, value, *rpc.Payload.Value)
}

func TestIncomingFindValueReturnsEmptyClosestContacts(t *testing.T) {
	node := Node{nil, Client{}, make(map[string]string)}
	c := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")
	node.RT = NewRoutingTable(c)
	network := InitServer(&node)

	payload := Payload{nil, nil, nil}
	senderID := "00000000000000000000000000000000FFFFFFFF"
	targetID := "00000000000000000000000000000000FFFFFFFF"
	rpc, _ := NewRPC(FindValue, senderID, targetID, payload)

	res, err := network.handleIncomingFindNodeRPC(rpc)
	assert.NoError(t, err)
	assert.Equal(t, []Contact(nil), res.Payload.Contacts)
}

func TestIncomingStoreSuccessfullyStoreValue(t *testing.T) {
	node := Node{nil, Client{}, make(map[string]string)}
	c := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")
	node.RT = NewRoutingTable(c)
	network := InitServer(&node)

	key := "hello"
	value := "good bye"
	payload := Payload{&key, &value, []Contact{}}

	rpc, _ := NewRPC(Store, "10000000000000000000000000000000FFFFFFFF", "00000000000000000000000000000000FFFFFFFF", payload)
	rpc, err := network.handleIncomingStoreRPC(rpc)

	val := node.searchLocalStore(key)

	assert.Nil(t, err)
	assert.NotNil(t, rpc)
	assert.Equal(t, value, *val)
}

func TestIncomingStoreBadInput(t *testing.T) {
	storeType := Store

	network := Server{}
	_, err := network.handleIncomingStoreRPC(nil)
	assert.Error(t, err)

	rpc := RPC{&storeType, nil, nil, nil, nil}
	_, err = network.handleIncomingStoreRPC(&rpc)
	assert.Error(t, err)

	payload := Payload{nil, nil, []Contact{}}
	rpc = RPC{&storeType, &payload, nil, nil, nil}
	_, err = network.handleIncomingStoreRPC(&rpc)
	assert.Error(t, err)
}

func TestNoTargetID(t *testing.T) {
	node := Node{}
	c := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")
	node.RT = NewRoutingTable(c)
	network := InitServer(&node)

	payload := Payload{nil, nil, []Contact{}}
	rpc, _ := NewRPC(Ping, "00000000000000000000000000000000FFFFFFFF", "", payload)
	rpc.TargetID = nil

	noTargetErr := errors.New(errNoTargetID)

	_, err1 := network.handleIncomingFindNodeRPC(rpc)
	_, err2 := network.handleIncomingFindValueRPC(rpc)

	assert.Equal(t, noTargetErr, err1)
	assert.Equal(t, noTargetErr, err2)
}

func TestPacketToIncomingChannel(t *testing.T) {
	node := Node{}
	server := InitServer(&node)
	go server.readIncomingChannel()
	ip := "127.0.0.1"

	payload := Payload{}
	rpc, _ := NewRPC(OK, "00000000000000000000000000000000FFFFFFFF", "00000000000000000000000000000000FFFFFFFF", payload)
	pkt := packet{rpc, ip, nil}
	server.incoming <- pkt

	val := <-server.outgoing

	assert.Equal(t, ip, val.ip)
	assert.Nil(t, val.rpc)
	assert.Nil(t, val.addr)
}

func TestHandleOutgoingChannel(t *testing.T) {
	node := Node{}
	server := InitServer(&node)
	addr, _ := net.ResolveUDPAddr(udpNetwork, "127.0.0.1:8080")
	server.conn, _ = net.ListenUDP(udpNetwork, addr)

	rpc, _ := NewRPC(Ping, "00000000000000000000000000000000FFFFFFFF", "00000000000000000000000000000000FFFFFFFF", Payload{nil, nil, []Contact{}})
	pkt := packet{rpc, "127.0.0.1", addr}
	server.outgoing <- pkt

	err := server.handleOutgoingChannel()
	assert.Nil(t, err)
}
