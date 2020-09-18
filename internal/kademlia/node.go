package kademlia

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Node struct {
	RT      *RoutingTable
	network Network
}

// InitNode initializes the Kademlia Node
// with a Routing Table and a Network
func (kademlia *Node) InitNode(id *NodeID) {
	kademlia.network = NewNetwork(kademlia)
	ip := kademlia.network.ip
	go kademlia.network.Listen(ip, "8080")

	me := NewContact(id, ip+":8080")
	rendezvousNode := NewContact(NewNodeID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.3:8080")
	kademlia.RT = NewRoutingTable(me)
	kademlia.RT.AddContact(rendezvousNode)
}

func (kademlia *Node) NodeLookup(target *Contact) {
	// TODO
}

func (kademlia *Node) FindValue(hash string) {
	fmt.Println("hash = ", hash)
	// TODO
}

func (kademlia *Node) StoreValue(data []byte) {
	fmt.Println("hash = ", data)
	// TODO
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
