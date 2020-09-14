package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewKademliaID(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewKademliaID("2111111400000000000000000000000000000000")

	assert.NotNil(t, id1)
	assert.Equal(t, id1.String(), "ffffffff00000000000000000000000000000000")
	assert.Equal(t, id2.String(), "2111111400000000000000000000000000000000")
}

func TestKademliaIDDistance(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewKademliaID("1111111100000000000000000000000000000000")
	id3 := NewKademliaID("2111111400000000000000000000000000000000")

	assert.Equal(t, id1.CalcDistance(id2).String(), "eeeeeeee00000000000000000000000000000000")
	assert.Equal(t, id3.CalcDistance(id1).String(), "deeeeeeb00000000000000000000000000000000")
	assert.Equal(t, id3.CalcDistance(id1), id1.CalcDistance(id3))
}

func TestKademliaIDEqual(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewKademliaID("1111111100000000000000000000000000000000")
	id3 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")

	assert.Equal(t, id1.Equals(id2), false)
	assert.Equal(t, id1.Equals(id3), true)
}

func TestKademliaIDLess(t *testing.T) {
	id1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewKademliaID("1111111100000000000000000000000000000000")
	id3 := NewKademliaID("2111111400000000000000000000000000000000")

	assert.Equal(t, id3.Less(id1), true)
	assert.Equal(t, id1.Less(id3), false)
	assert.Equal(t, id1.Less(id2), false)
}
