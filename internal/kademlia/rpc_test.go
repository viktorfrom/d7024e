package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRPCMarshal(t *testing.T) {

}

func TestRPCUnmarshal(t *testing.T) {
	msg := "hello"
	originalRPC := NewRPC(Ping, []byte(msg))

	data, _ := MarshalRPC(originalRPC)
	marshalledRPC, _ := UnmarshalRPC(data)

	assert.Equal(t, msg, *marshalledRPC.Payload)
	assert.Equal(t, Ping, *marshalledRPC.Type)
}

func TestRPCCheckSameID(t *testing.T) {
	msg := "hello"
	originalRPC := NewRPC(Store, []byte(msg))
	originalID := *originalRPC.ID

	data, _ := MarshalRPC(originalRPC)
	marshalledRPC, _ := UnmarshalRPC(data)
	newID := *marshalledRPC.ID

	assert.Equal(t, originalID, newID)
}
