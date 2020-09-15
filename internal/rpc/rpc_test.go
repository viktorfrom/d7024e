package rpc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRPCUnmarshal(t *testing.T) {
	msg := "hello"
	originalRPC, _ := NewRPC(Ping, []byte(msg))

	data, _ := MarshalRPC(*originalRPC)
	marshalledRPC, _ := UnmarshalRPC(data)

	assert.Equal(t, msg, *marshalledRPC.Payload)
	assert.Equal(t, Ping, *marshalledRPC.Type)
}

func TestRPCWrongDataUnmarshal(t *testing.T) {
	data := []byte{1, 0, 1, 0, 1}
	_, err := UnmarshalRPC(data)
	assert.Error(t, err)
}

func TestRPCValidateID(t *testing.T) {
	msg := []byte("hello")
	originalRPC, _ := NewRPC(Store, msg)
	originalID := *originalRPC.ID

	data, _ := MarshalRPC(*originalRPC)
	marshalledRPC, _ := UnmarshalRPC(data)
	newID := *marshalledRPC.ID

	assert.Equal(t, originalID, newID)
}

func TestEmptyRPC(t *testing.T) {
	originalRPC := RPC{}

	data, err := MarshalRPC(originalRPC)
	assert.Nil(t, err)

	returnedRPC, err := UnmarshalRPC(data)

	assert.Nil(t, err)
	assert.Equal(t, &originalRPC, returnedRPC)
}

func TestNewRPCCorrectTypes(t *testing.T) {
	msg := []byte("good bye")

	for _, rpcType := range rpcTypes {
		_, err := NewRPC(rpcType, msg)
		assert.NoError(t, err)
	}
}

func TestNewRPCWrongType(t *testing.T) {
	msg := []byte("good bye")
	_, err := NewRPC("wrong type", msg)
	assert.Error(t, err)
}
