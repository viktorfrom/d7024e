package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func fillBucket(b *bucket) *bucket {
	for i := 0; i < bucketSize*2; i++ {
		contact := NewContact(NewRandomKademliaID(), "10.0.8.2")
		b.AddContact(contact)
	}
	return b
}

func TestNewBucket(t *testing.T) {
	bucket1 := newBucket()
	bucket2 := newBucket()

	assert.NotNil(t, bucket1)
	assert.NotNil(t, bucket2)
}

func TestBucketAddContact(t *testing.T) {
	bucket1 := newBucket()

	kID := NewRandomKademliaID()
	contact1 := NewContact(kID, "10.0.8.2")

	assert.Nil(t, bucket1.list.Front())
	bucket1.AddContact(contact1)
	assert.Equal(t, contact1, bucket1.list.Front().Value)

	bucket1 = fillBucket(bucket1)

	// check so that contact1 gets moved from the back of the list to the front
	assert.NotEqual(t, contact1, bucket1.list.Front().Value)
	assert.Equal(t, contact1, bucket1.list.Back().Value)
	bucket1.AddContact(contact1)
	assert.Equal(t, contact1, bucket1.list.Front().Value)
}

func TestBucketGetContact(t *testing.T) {
	bucket1 := newBucket()
	bucket1 = fillBucket(bucket1)
	contact1 := bucket1.list.Front().Value.(Contact)

	kID := NewRandomKademliaID()

	contact1.CalcDistance(kID)
	assert.Equal(t, contact1, bucket1.GetContactAndCalcDistance(kID)[0])
	assert.NotNil(t, bucket1.GetContactAndCalcDistance(kID)[1].distance)
	assert.NotNil(t, bucket1.GetContactAndCalcDistance(kID)[4].distance)
}

func TestBucketLen(t *testing.T) {
	bucket1 := newBucket()
	contact1 := NewContact(NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "10.0.8.2")

	assert.Equal(t, 0, bucket1.Len())
	bucket1.AddContact(contact1)
	assert.Equal(t, 1, bucket1.Len())

	bucket1 = fillBucket(bucket1)
	assert.Equal(t, bucketSize, bucket1.Len())
}
