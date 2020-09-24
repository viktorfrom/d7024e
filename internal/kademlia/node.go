package kademlia

import (
	"crypto/sha1"
	"fmt"
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
	kademlia.RT = NewRoutingTable(me)

	rendezvousNode := NewContact(rendezvousID, "10.0.8.3:8080")
	kademlia.JoinNetwork(rendezvousNode)

	kademlia.content = make(map[string]string)
}

func (kademlia *Node) NodeLookup(target *Contact) {

	// TODO: support for parallelism alpha = ~3
	// TODO: if a closer node is not found, the initiating node sends a FIND_* RPC to each of the
	// k closest nodes that it has not already queried.
	closestsContacts := kademlia.RT.FindClosestContacts(target.ID, BucketSize)
	shortList := ContactCandidates{closestsContacts}

	closestNode := shortList.contacts[0]
	// fmt.Println("table = ", closestNode)

	for {
		fmt.Println("table = ", shortList)

		if closestNode.ID.Equals(target.ID) {
			fmt.Println("node found = ", closestNode)
			break

		} else {

			rpc, err := kademlia.network.SendFindContactMessage(&closestNode, &kademlia.RT.me)

			// remove current/first node from shortlist
			if len(shortList.contacts) > 0 {
				shortList.contacts = shortList.contacts[1:]
			}

			// append contacts to shortlist if err is none
			if err == nil {
				for i := 0; i < len(rpc.Payload.Contacts); i++ {
					shortList.contacts = appendUnique(shortList.contacts, rpc.Payload.Contacts[i])
				}
			}

			shortList.Sort()

			// update closest node if first element distance is shorter
			if len(shortList.contacts) > 0 || shortList.contacts[0].Less(target) {
				closestNode = shortList.contacts[0]
			}

			// sleep for testing
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func appendUnique(slice []Contact, i Contact) []Contact {
	for _, ele := range slice {
		if ele == i {
			return slice
		}
	}

	return append([]Contact{i}, slice...)
}

func (kademlia *Node) FindValue(hash string) {
	sha1 := sha1.Sum([]byte(hash))
	var content = kademlia.content[string(sha1[:])]
	if content == "" {
		fmt.Println("Content not found!")
	} else {
		// return content
		fmt.Println("Content = ", content)
	}
	// return content
}

func (kademlia *Node) StoreValue(data string) {
	sha1 := sha1.Sum([]byte(data))
	kademlia.content[string(sha1[:])] = data
}

func (kademlia *Node) Ping() {
	target := &kademlia.RT.FindClosestContacts(kademlia.RT.me.ID, BucketSize)[0]
	rpc, err := kademlia.network.SendPingMessage(target, &kademlia.RT.me)

	if err != nil {
		log.Warn(err)
		kademlia.RT.RemoveContact(*target)
	} else if *rpc.Type == "OK" {
		kademlia.RT.AddContact(*target)
	}
}

// SearchStore looks for a value in the node's store. Returns the value
// if found else nil.
func (kademlia *Node) SearchStore(key string) *string {
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
		kademlia.NodeLookup(&contact)
	}
}

// JoinNetwork add a target node to the routing table, do a Node Lookup on
// the current node (not the target) and then refresh all buckets
func (kademlia *Node) JoinNetwork(target Contact) {

	kademlia.RT.AddContact(target)

	// kademlia.NodeLookup(kademlia.RT.GetMe())

	// kademlia.refreshNodes()
}
