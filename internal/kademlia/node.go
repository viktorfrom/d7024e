package kademlia

import (
	"crypto/sha1"
	"encoding/hex"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

type Node struct {
	RT      *RoutingTable
	network Network
	content map[string]string
}

// InitNode initializes the Kademlia Node
// with a Routing Table and a Network
func (kademlia *Node) InitNode() {
	kademlia.network = NewNetwork(kademlia)
	ip := kademlia.network.ip

	var id *NodeID

	rendezvousID := NewNodeID("00000000000000000000000000000000FFFFFFFF")

	// set a specific ID to the rendezvous node, the node that has the address "10.0.8.3"
	if ip == "10.0.8.3" {
		id = rendezvousID
	} else {
		// for all nodes that is not the rendezvous node set a random ID
		id = NewRandomNodeID()
	}

	go kademlia.network.Listen(ip, "8080")

	me := NewContact(id, ip+":8080")
	me.CalcDistance(me.ID)
	kademlia.RT = NewRoutingTable(me)

	if ip != "10.0.8.3" {
		rendezvousNode := NewContact(rendezvousID, "10.0.8.3:8080")

		// wait a second before trying to join the network to allow the rendezvousNode
		// to become active
		time.Sleep(1 * time.Second)
		kademlia.JoinNetwork(rendezvousNode)
	}

	kademlia.content = make(map[string]string)
}

func (kademlia *Node) NodeLookup(targetID *NodeID) []Contact {
	alpha := 1
	shortList := ContactCandidates{kademlia.RT.FindClosestContacts(targetID, alpha)}

	// set a temporary value to currentClosest that is the furthest away a node can be
	currentClosest := NewContact(NewNodeID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "")
	currentClosest.distance = NewNodeID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

	// a list of nodes to know which nodes has been probed already
	probedNodes := ContactCandidates{}

	for {
		updateClosest := false
		numProbed := 0

		for i := 0; i < shortList.Len() && numProbed < alpha; i++ {

			if probedNodes.Contains(shortList.contacts[i]) {
				continue
			} else {
				rpc, err := kademlia.network.SendFindContactMessage(&shortList.contacts[i], &kademlia.RT.me)

				// if a node responds with an error remove that node from the shortlist and from the bucket
				if err != nil {
					log.Warn(err)
					kademlia.RT.RemoveContact(shortList.contacts[i])
					shortList.contacts = append(shortList.contacts[:i], shortList.contacts[i+1:]...)
					continue
				}

				probedNodes.Append([]Contact{shortList.contacts[i]})

				bucket := kademlia.RT.buckets[kademlia.RT.getBucketIndex(shortList.contacts[i].ID)]

				kademlia.updateBucket(*bucket, shortList.contacts[i])

				// append contacts to shortlist if err is none
				for i := 0; i < len(rpc.Payload.Contacts); i++ {
					rpc.Payload.Contacts[i].CalcDistance(targetID)
				}

				// if the closest node in the payload is less than the currentClosest
				// update the shortlist and the currentClosest node
				if rpc.Payload.Contacts[0].Less(&currentClosest) {
					currentClosest = rpc.Payload.Contacts[0]
					shortList.AppendUnique(rpc.Payload.Contacts)
					shortList.Sort()
					if shortList.Len() >= BucketSize {
						shortList.contacts = shortList.contacts[:BucketSize]
					}

					updateClosest = true
				}

				numProbed++
			}
		}

		if !updateClosest || probedNodes.Len() >= BucketSize {
			break
		}
	}

	return shortList.contacts
}

func (kademlia *Node) FindValue(hash string) string {

	if content, ok := kademlia.content[hash]; ok {
		return content

	} else {
		alpha := 1
		shortList := ContactCandidates{kademlia.RT.FindClosestContacts(NewNodeID(hash), alpha)}

		// set a temporary value to currentClosest that is the furthest away a node can be
		currentClosest := NewContact(NewNodeID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF"), "")
		currentClosest.distance = NewNodeID("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

		// a list of nodes to know which nodes has been probed already
		probedNodes := ContactCandidates{}

		for {
			updateClosest := false
			numProbed := 0

			for i := 0; i < shortList.Len() && numProbed < alpha; i++ {

				if probedNodes.Contains(shortList.contacts[i]) {
					continue
				} else {
					rpc, err := kademlia.network.SendFindDataMessage(&shortList.contacts[i], &kademlia.RT.me, hash)

					if *rpc.Payload.Value != "" {
						return *rpc.Payload.Value
					}

					// if a node responds with an error remove that node from the shortlist and from the bucket
					if err != nil {
						log.Warn(err)
						kademlia.RT.RemoveContact(shortList.contacts[i])
						shortList.contacts = append(shortList.contacts[:i], shortList.contacts[i+1:]...)
						continue
					}

					probedNodes.Append([]Contact{shortList.contacts[i]})

					bucket := kademlia.RT.buckets[kademlia.RT.getBucketIndex(shortList.contacts[i].ID)]

					// if there is space in the bucket add the node
					if bucket.Len() < BucketSize {
						kademlia.RT.AddContact(shortList.contacts[i])
					} else {
						// if there is no space in the bucket ping the least recently seen node
						kademlia.Ping(bucket.GetFirst())

						// if there now is space in the bucket add the node
						if bucket.Len() < BucketSize {
							kademlia.RT.AddContact(shortList.contacts[i])
						}
					}

					// append contacts to shortlist if err is none
					for i := 0; i < len(rpc.Payload.Contacts); i++ {
						rpc.Payload.Contacts[i].CalcDistance(NewNodeID(hash))
					}

					// if the closest node in the payload is less than the currentClosest
					// update the shortlist and the currentClosest node
					if rpc.Payload.Contacts[0].Less(&currentClosest) {
						currentClosest = rpc.Payload.Contacts[0]
						shortList.AppendUnique(rpc.Payload.Contacts)
						shortList.Sort()
						if shortList.Len() >= BucketSize {
							shortList.contacts = shortList.contacts[:BucketSize]
						}

						updateClosest = true
					}

					numProbed++
				}
			}

			if !updateClosest || probedNodes.Len() >= BucketSize {
				break
			}
		}

		// return shortList.contacts

	}

	return "No value found!"
}

// StoreValue takes some data, hashes it with SHA1 and finds the k closest
// nodes to that hash, then sends a store RPC to those k nodes
func (kademlia *Node) StoreValue(data string) {
	sha1 := sha1.Sum([]byte(data))
	key := hex.EncodeToString(sha1[:])

	// find the K closest nodes to the hashed value in the whole Kademlia network
	targetID := NewNodeID(key)
	nodes := kademlia.NodeLookup(targetID)

	// for each of the closest nodes send a store RPC
	for _, node := range nodes {
		_, err := kademlia.network.SendStoreMessage(&node, &kademlia.RT.me, key, data)

		if err != nil {
			log.Warn(err)
			kademlia.RT.RemoveContact(node)
		} else {
			bucket := kademlia.RT.buckets[kademlia.RT.getBucketIndex(node.ID)]

			kademlia.updateBucket(*bucket, node)
		}
	}
}

// Ping sends a ping message to a target node
// if the node responds move it to the end of the bucket it exists in
// if the node does not respond remove it from the bucket
func (kademlia *Node) Ping(target *Contact) {
	rpc, err := kademlia.network.SendPingMessage(target, &kademlia.RT.me)

	if err != nil {
		log.Warn(err)
		kademlia.RT.RemoveContact(*target)
	} else if *rpc.Type == "OK" {
		kademlia.RT.AddContact(*target)
	}
}

// updateBucket checks if a contact should be added to a bucket if it does not exist,
// removes a stale first node in the bucket and replace it with the new node
// or a active old node from the front to the back
func (kademlia *Node) updateBucket(bucket bucket, contact Contact) {
	// if there is space in the bucket add the node
	if bucket.Len() < BucketSize || bucket.Contains(contact) {
		kademlia.RT.AddContact(contact)
	} else {
		// if there is no space in the bucket ping the least recently seen node
		kademlia.Ping(bucket.GetFirst())

		// if there now is space in the bucket add the node
		if bucket.Len() < BucketSize {
			kademlia.RT.AddContact(contact)
		}
	}
}

// searchLocalStore looks for a value in the node's store. Returns the value
// if found else nil.
func (kademlia *Node) searchLocalStore(key string) *string {
	value, exists := kademlia.content[key]
	if exists {
		return nil
	}
	return &value
}

// generate a random ID that is inside a given bucket
func generateRefreshNodeValue(bucketIndex int, seed int64) *NodeID {
	bytePos := 19 - (bucketIndex / 8) // position of the highest byte of the ID
	offset := bucketIndex % 8

	nodeValue := NewNodeID("0000000000000000000000000000000000000000")

	t := 0
	t = 1 << offset

	nodeValue[bytePos] = byte(t)
	rand.Seed(int64(seed))

	// generate a random byte for each byte position from the end of the string to the bytePos
	for i := 19; i > bytePos; i-- {
		scew := uint8(rand.Intn(bucketIndex))
		nodeValue[i] ^= byte(scew)
	}

	return nodeValue
}

func (kademlia *Node) refreshNodes() {
	for i := 1; i > 159; i++ {
		nodeID := generateRefreshNodeValue(i, time.Now().UTC().UnixNano())
		contact := NewContact(nodeID, "")
		kademlia.NodeLookup(contact.ID)
	}
}

// JoinNetwork add a target node to the routing table, do a Node Lookup on
// the current node (not the target) and then refresh all buckets
func (kademlia *Node) JoinNetwork(target Contact) {

	kademlia.RT.AddContact(target)

	kademlia.NodeLookup(kademlia.RT.GetMe().ID)

	kademlia.refreshNodes()
}

func (kademlia *Node) insertLocalStore(key string, value string) {
	kademlia.content[key] = value
}
