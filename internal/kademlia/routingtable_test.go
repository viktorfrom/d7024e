package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewRoutingTable(t *testing.T) {
	me := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000")
	rt := NewRoutingTable(me)
	assert.NotNil(t, rt)
	assert.Equal(t, me, rt.me)

	// check that the correct amount of buckets are initialized
	assert.Equal(t, IDLength*8, len(rt.buckets))
	for i := 0; i < IDLength*8; i++ {
		assert.NotNil(t, rt.buckets[i])
	}
}

func TestRoutingTable(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	c1 := NewContact(NewKademliaID("1111111400000000000000000000000000000000"), "localhost:8002")
	c2 := NewContact(NewKademliaID("2111111400000000000000000000000000000000"), "localhost:8002")
	rt.AddContact(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8001"))
	rt.AddContact(NewContact(NewKademliaID("1111111100000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(c1)
	rt.AddContact(c2)

	contacts := rt.FindClosestContacts(NewKademliaID("2111111400000000000000000000000000000000"), 2)

	assert.Equal(t, c2.ID, contacts[0].ID)
	assert.Equal(t, c1.ID, contacts[1].ID)
}

func TestGetBucketIndex(t *testing.T) {
	rt := NewRoutingTable(NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "localhost:8000"))

	kID1 := NewKademliaID("FFFFFFFF00000000000000000000000000000000")
	kID2 := NewKademliaID("1111111100000000000000000000000000000000")
	kID3 := NewKademliaID("1111111400000000000000000000000000000000")
	kID4 := NewKademliaID("21111114000000000000000000000000FFFFFFFF")
	rt.AddContact(NewContact(kID1, "localhost:8001"))
	rt.AddContact(NewContact(kID2, "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111200000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(NewKademliaID("1111111300000000000000000000000000000000"), "localhost:8002"))
	rt.AddContact(NewContact(kID3, "localhost:8002"))
	rt.AddContact(NewContact(kID4, "localhost:8002"))

	assert.Equal(t, rt.getBucketIndex(kID2), rt.getBucketIndex(kID3))
	assert.Equal(t, 159, rt.getBucketIndex(kID1))
}
