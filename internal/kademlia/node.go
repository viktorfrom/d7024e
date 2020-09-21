package kademlia

import (
	"crypto/md5"
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
	c1 := NewContact(NewNodeID("1111111400000000000000000000000000000000"), "localhost:8002")
	c2 := NewContact(NewNodeID("2111111400000000000000000000000000000000"), "localhost:8002")
	c3 := NewContact(NewNodeID("3111111400000000000000000000000000000000"), "localhost:8002")
	c4 := NewContact(NewNodeID("4111111400000000000000000000000000000000"), "localhost:8002")

	kademlia.RT.AddContact(c4)
	kademlia.RT.AddContact(c1)
	kademlia.RT.AddContact(c2)
	kademlia.RT.AddContact(c3)

	table := kademlia.RT.FindClosestContacts(target.ID, BucketSize)

	for i := 0; i < len(table); i++ {
		// fmt.Println("table = ", table[i], "target = ", target.ID)
		if table[i].ID.Equals(target.ID) {
			fmt.Println("node found = ", table[i])
		} else {
			// TODO: add iterative/recursive RPC call
		}
	}
}

func (kademlia *Node) FindValue(hash string) string {
	md5 := md5.Sum([]byte(hash))
	var content = kademlia.content[string(md5[:])]
	if content == "" {
		fmt.Println("Content not found!")
		// } else {
		// 	return content
		// fmt.Println("Content = ", content)
	}
	return content
}

func (kademlia *Node) StoreValue(data string) {
	md5 := md5.Sum([]byte(data))
	kademlia.content[string(md5[:])] = data
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

	kademlia.NodeLookup(kademlia.RT.GetMe())

	kademlia.refreshNodes()
}
