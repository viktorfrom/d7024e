package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoutingTable(t *testing.T) {
	me := NewContact(NewNodeID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	rt := NewRoutingTable(me)
	assert.NotNil(t, rt)
	assert.Equal(t, me, rt.me)

	// check that the correct amount of buckets are initialized
	assert.Equal(t, IDLength*8, len(rt.buckets))
	for i := 0; i < IDLength*8; i++ {
		assert.NotNil(t, rt.buckets[i])
	}
}

func TestSmallRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewNodeID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	c1 := NewContact(NewNodeID("1111111400000000000000000000000000000000"), "localhost:8002")
	c2 := NewContact(NewNodeID("2111111400000000000000000000000000000000"), "localhost:8002")
	rt.AddContact(NewContact(NewNodeID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewNodeID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewNodeID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(c1)
	rt.AddContact(c2)

	contacts := rt.FindClosestContacts(NewNodeID("2111111400000000000000000000000000000000"), 2)

	assert.Equal(t, c2.ID, contacts[0].ID)
	assert.Equal(t, c1.ID, contacts[1].ID)

	rt.AddContact(NewContact(NewRandomNodeID(), "localhost:8002"))
	rt.AddContact(NewContact(NewRandomNodeID(), "localhost:8002"))

	contacts = rt.FindClosestContacts(NewRandomNodeID(), 10)
	assert.NotNil(t, contacts)
}

func TestBigRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewNodeID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	for i := 0; i < 1000; i++ {
		rt.AddContact(NewContact(NewRandomNodeID(), "localhost:8001"))
	}

	contacts := rt.FindClosestContacts(NewNodeID("2111111400000000000000000000000000000000"), 20)
	assert.NotNil(t, contacts)
	contacts = rt.FindClosestContacts(NewNodeID("2111111400000000000000000000000000000000"), 200)
	assert.NotNil(t, contacts)
}

func TestGetBucketIndex(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewNodeID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	kID1 := NewNodeID("FFFFFFFF00000000000000000000000000000000")
	kID2 := NewNodeID("1111111100000000000000000000000000000000")
	kID3 := NewNodeID("1111111400000000000000000000000000000000")
	kID4 := NewNodeID("21111114000000000000000000000000FFFFFFFF")
	rt.AddContact(NewContact(kID1, "localhost:8001"))
	rt.AddContact(NewContact(kID2, "localhost:8002"))
	rt.AddContact(NewContact(NewNodeID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewNodeID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(kID3, "localhost:8002"))
	rt.AddContact(NewContact(kID4, "localhost:8002"))

	assert.Equal(t, rt.getBucketIndex(kID2), rt.getBucketIndex(kID3))
	assert.Equal(t, 159, rt.getBucketIndex(kID1))
}

func TestRemoveFromRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewNodeID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	c1 := NewContact(NewNodeID("1111111400000000000000000000000000000000"), "localhost:8002")
	c2 := NewContact(NewNodeID("2111111400000000000000000000000000000000"), "localhost:8002")
	c3 := NewContact(NewNodeID("1111111100000000000000000000000000000000"), "localhost:8002")

	rt.AddContact(NewContact(NewNodeID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewNodeID("ffff111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(c1)
	rt.AddContact(c2)
	rt.AddContact(c3)

	contacts := rt.FindClosestContacts(NewNodeID("2111111400000000000000000000000000000000"), 2)
	assert.Equal(t, c2.ID, contacts[0].ID)
	assert.Equal(t, c1.ID, contacts[1].ID)

	rt.RemoveContact(c1)
	contacts = rt.FindClosestContacts(NewNodeID("2111111400000000000000000000000000000000"), 2)

	assert.Equal(t, c2.ID, contacts[0].ID)
	assert.Equal(t, c3.ID, contacts[1].ID)
}

func TestGetMe(t *testing.T) {
	c1 := NewContact(NewNodeID("1111111400000000000000000000000000000000"), "localhost:8002")
	rt := NewRoutingTable(c1)

	assert.Equal(t, c1, *rt.GetMe())
	assert.Equal(t, c1.ID, rt.GetMeID())
}
