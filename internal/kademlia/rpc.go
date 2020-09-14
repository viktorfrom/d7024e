package network

import (
	"encoding/json"

	"github.com/viktorfrom/d7024e-kademlia/pkg/randarr"
)

type RPCType string

// RPC type definitions
const (
	Ping      = RPCType("PING")
	Store     = RPCType("STORE")
	FindValue = RPCType("FIND_VALUE")
	FindNode  = RPCType("FIND_NODE")
	OK        = RPCType("OK")
)

// RPC ...
type RPC struct {
	Type    *RPCType `json:"type"`
	Payload *string  `json:"payload"`
	ID      *string  `json:"id"`
}

// NewRPC creates a new RPC with a random ID added to it
func NewRPC(rpc RPCType, data []byte) RPC {
	payload := string(data)
	randomStr := randarr.RandomHexString(20)
	randomID := string(randomStr)

	return RPC{&rpc, &payload, &randomID}
}

// MarshalRPC serializes the RPC struct and returns the result as a byte array
func MarshalRPC(rpc RPC) ([]byte, error) {
	var data []byte
	data, err := json.Marshal(rpc)
	if err != nil {
		return []byte{}, err
	}

	return data, nil
}

// UnmarshalRPC deserializes the given byte array and returns an RPC
func UnmarshalRPC(data []byte) (*RPC, error) {
	rpc := RPC{}
	err := json.Unmarshal(data, &rpc)
	if err != nil {
		return nil, err
	}

	return &rpc, nil
}
