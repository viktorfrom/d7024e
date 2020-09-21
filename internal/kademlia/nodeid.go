package kademlia

import (
	"encoding/hex"
	"math/rand"
	"time"
)

// IDLength the static number of bytes in a NodeID
const IDLength = 20

// NodeID type definition of a NodeID
type NodeID [IDLength]byte

// NewNodeID returns a new instance of a NodeID based on the string input
func NewNodeID(data string) *NodeID {
	decoded, _ := hex.DecodeString(data)

	newNodeID := NodeID{}
	for i := 0; i < IDLength; i++ {
		newNodeID[i] = decoded[i]
	}

	return &newNodeID
}

// NewRandomNodeID returns a new instance of a random NodeID,
// change this to a better version if you like
func NewRandomNodeID() *NodeID {
	rand.Seed(time.Now().UTC().UnixNano())
	newNodeID := NodeID{}
	for i := 0; i < IDLength; i++ {
		newNodeID[i] = uint8(rand.Intn(256))
	}
	return &newNodeID
}

// Less returns true if NodeID < otherNodeID (bitwise)
func (nodeID NodeID) Less(otherNodeID *NodeID) bool {
	for i := 0; i < IDLength; i++ {
		if nodeID[i] != otherNodeID[i] {
			return nodeID[i] < otherNodeID[i]
		}
	}
	return false
}

// Equals returns true if NodeID == otherNodeID (bitwise)
func (nodeID NodeID) Equals(otherNodeID *NodeID) bool {
	for i := 0; i < IDLength; i++ {
		if nodeID[i] != otherNodeID[i] {
			return false
		}
	}
	return true
}

// CalcDistance returns a new instance of a NodeID that is built
// through a bitwise XOR operation betweeen NodeID and target
func (nodeID NodeID) CalcDistance(target *NodeID) *NodeID {
	result := NodeID{}
	for i := 0; i < IDLength; i++ {
		result[i] = nodeID[i] ^ target[i]
	}
	return &result
}

// String returns a simple string representation of a NodeID
func (nodeID *NodeID) String() string {
	return hex.EncodeToString(nodeID[0:IDLength])
}
