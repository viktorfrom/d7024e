package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRPCUnmarshal(t *testing.T) {
	msg := "hello"
	contact := NewContact(NewRandomKademliaID(), "10.0.8.2")
	payload := Payload{&msg, []Contact{contact}}
	originalRPC, _ := NewRPC(Ping, "10.0.8.1", payload)

	data, _ := MarshalRPC(*originalRPC)
	marshalledRPC, _ := UnmarshalRPC(data)

	assert.Equal(t, msg, *marshalledRPC.Payload.Value)
	assert.Equal(t, contact, marshalledRPC.Payload.Contacts[0])
	assert.Equal(t, Ping, *marshalledRPC.Type)
}

func TestRPCWrongDataUnmarshal(t *testing.T) {
	data := []byte{1, 0, 1, 0, 1}
	_, err := UnmarshalRPC(data)
	assert.Error(t, err)
}

func TestRPCValidateID(t *testing.T) {
	msg := "hello"
	payload := Payload{&msg, nil}
	originalRPC, _ := NewRPC(Store, "10.0.8.1", payload)
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
	msg := "good bye"
	payload := Payload{&msg, nil}

	for _, rpcType := range rpcTypes {
		_, err := NewRPC(rpcType, "10.0.8.1", payload)
		assert.NoError(t, err)
	}
}

func TestNewRPCWrongType(t *testing.T) {
	msg := "good bye"
	payload := Payload{&msg, nil}

	_, err := NewRPC("wrong type", "10.0.8.1", payload)
	assert.Error(t, err)
}
