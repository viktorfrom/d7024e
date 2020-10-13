package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitClient(t *testing.T) {
	client := InitClient()
	assert.NotNil(t, client.ip)
}

func TestGetLocalIp(t *testing.T) {
	client := Client{}
	ip := client.GetLocalIP()
	assert.NotEqual(t, "", ip)
	assert.Equal(t, ip, client.GetLocalIP())
}

// func TestSendPingError(t *testing.T) {
// 	node := Node{}
// 	c := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")
// 	node.RT = NewRoutingTable(c)
// 	network := NewNetwork(&node)

// 	_, err := network.SendPingMessage(nil, &c)
// 	assert.Error(t, err)
// }

// func TestSendFindNodeError(t *testing.T) {
// 	network := Network{}
// 	_, err := network.SendFindContactMessage(nil, nil, nil)
// 	assert.Error(t, err)
// }

// func TestSendFindValueError(t *testing.T) {
// 	node := Node{}
// 	c := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.1:8080")
// 	node.RT = NewRoutingTable(c)
// 	network := NewNetwork(&node)

// 	_, err := network.SendFindDataMessage(&c, nil, "testkey")
// 	assert.Error(t, err)
// }

// func TestSendStoreError(t *testing.T) {
// 	network := Client{}
// 	_, err := network.SendStoreMessage(nil, nil, "key", "value")
// 	assert.Error(t, err)
// }

// func TestSendRPCNoNetworkAvailable(t *testing.T) {
// 	timeout = 0 * time.Second

// 	network := Client{}
// 	nodeID := NewNodeID("00000000000000000000000000000000FFFFFFFF")
// 	c := NewContact(nodeID, "10.0.8.1:8080")

// 	payload := Payload{nil, nil, []Contact{}}
// 	_, err := network.sendRPC(&c, Ping, nodeID, nodeID, payload)
// 	assert.Error(t, err)
// }
