package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewContact(t *testing.T) {
	contact1 := NewContact(NewNodeID("ffffffff00000000000000000000000000000000"), "10.0.8.0")
	contact2 := NewContact(NewNodeID("00000000000000000000000000000000ffffffff"), "10.0.8.5")

	assert.NotNil(t, contact1)
	assert.Equal(t, "contact(\"ffffffff00000000000000000000000000000000\", \"10.0.8.0\")", contact1.String())
	assert.Equal(t, "contact(\"00000000000000000000000000000000ffffffff\", \"10.0.8.5\")", contact2.String())
}

func TestContactDistance(t *testing.T) {
	contact1 := NewContact(NewNodeID("FFFFFFFF00000000000000000000000000000000"), "10.0.8.0")
	contact2 := NewContact(NewNodeID("1111111100000000000000000000000000000000"), "10.0.8.5")

	contact1.CalcDistance(contact2.ID)
	contact2.CalcDistance(contact1.ID)
	assert.Equal(t, "eeeeeeee00000000000000000000000000000000", contact1.distance.String())
	assert.Equal(t, contact2.distance.String(), contact1.distance.String())
}

func TestContactLess(t *testing.T) {
	contact1 := NewContact(NewNodeID("FFFFFFFF00000000000000000000000000000000"), "10.0.8.0")
	contact2 := NewContact(NewNodeID("1111111100000000000000000000000000000000"), "10.0.8.5")
	contact3 := NewContact(NewNodeID("0000000000000000000000000000000011111111"), "10.0.8.6")

	contact1.CalcDistance(contact2.ID)
	contact2.CalcDistance(contact3.ID)

	assert.Equal(t, false, contact1.Less(&contact2))
	assert.Equal(t, true, contact2.Less(&contact1))
}

func TestContactCandidatesAppend(t *testing.T) {
	contactCandidates := ContactCandidates{}
	assert.Equal(t, 0, contactCandidates.Len())

	contacts := []Contact{NewContact(NewNodeID("FFFFFFFF00000000000000000000000000000000"), "10.0.8.0")}
	contactCandidates.Append(contacts)
	assert.Equal(t, 1, contactCandidates.Len())

	contacts = []Contact{NewContact(NewNodeID("FFFFFFFF00000000000000000000000000000000"), "10.0.8.0"), NewContact(NewNodeID("FFFFFFFF00000000000000000000000000000000"), "10.0.8.0")}
	contactCandidates.Append(contacts)
	assert.Equal(t, 3, contactCandidates.Len())
}

func TestContactCandidatesGet(t *testing.T) {
	contact1 := NewContact(NewNodeID("FFFFFFFF00000000000000000000000000000000"), "10.0.8.0")
	contact2 := NewContact(NewNodeID("00000000FFFFFFFF000000000000000000000000"), "10.0.8.0")
	contact3 := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.0")

	contactCandidates := ContactCandidates{}
	assert.Equal(t, []Contact(nil), contactCandidates.GetContacts(0))

	contacts := []Contact{contact1, contact2, contact3}
	contactCandidates.Append(contacts)
	assert.Equal(t, []Contact{contact1}, contactCandidates.GetContacts(1))
}

func TestContactCandidatesSwap(t *testing.T) {
	contact1 := NewContact(NewNodeID("FFFFFFFF00000000000000000000000000000000"), "10.0.8.2")
	contact2 := NewContact(NewNodeID("00000000FFFFFFFF000000000000000000000000"), "10.0.8.3")
	contact3 := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.4")

	contacts := []Contact{contact1, contact2, contact3}
	contactCandidates := ContactCandidates{contacts}
	assert.Equal(t, []Contact{contact1}, contactCandidates.GetContacts(1))

	contactCandidates.Swap(0, 2)
	assert.Equal(t, []Contact{contact3}, contactCandidates.GetContacts(1))

	contactCandidates.Swap(0, 0)
	assert.Equal(t, []Contact{contact3}, contactCandidates.GetContacts(1))
}
