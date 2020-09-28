package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewNodeID(t *testing.T) {
	id1 := NewNodeID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewNodeID("2111111400000000000000000000000000000000")

	assert.NotNil(t, id1)
	assert.Equal(t, id1.String(), "ffffffff00000000000000000000000000000000")
	assert.Equal(t, id2.String(), "2111111400000000000000000000000000000000")
}

func TestNewRandomNodeID(t *testing.T) {
	id1 := NewRandomNodeID()

	assert.NotNil(t, id1)
}

func TestNodeIDDistance(t *testing.T) {
	id1 := NewNodeID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewNodeID("1111111100000000000000000000000000000000")
	id3 := NewNodeID("2111111400000000000000000000000000000000")

	assert.Equal(t, id1.CalcDistance(id2).String(), "eeeeeeee00000000000000000000000000000000")
	assert.Equal(t, id3.CalcDistance(id1).String(), "deeeeeeb00000000000000000000000000000000")
	assert.Equal(t, id3.CalcDistance(id1), id1.CalcDistance(id3))
}

func TestNodeIDEqual(t *testing.T) {
	id1 := NewNodeID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewNodeID("1111111100000000000000000000000000000000")
	id3 := NewNodeID("FFFFFFFF00000000000000000000000000000000")

	assert.Equal(t, id1.Equals(id2), false)
	assert.Equal(t, id1.Equals(id3), true)
}

func TestNodeIDLess(t *testing.T) {
	id1 := NewNodeID("FFFFFFFF00000000000000000000000000000000")
	id2 := NewNodeID("1111111100000000000000000000000000000000")
	id3 := NewNodeID("2111111400000000000000000000000000000000")

	assert.Equal(t, id3.Less(id1), true)
	assert.Equal(t, id1.Less(id3), false)
	assert.Equal(t, id1.Less(id2), false)
	assert.Equal(t, id1.Less(id1), false)
}
