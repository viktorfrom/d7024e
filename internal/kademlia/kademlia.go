package kademlia

import "fmt"

type Kademlia struct {
	RT      *RoutingTable
	network Network
}

// InitNode initializes the Kademlia Node
// with a Routing Table and a Network
func (kademlia *Kademlia) InitNode(id *KademliaID) {
	kademlia.network = Network{}
	ip := kademlia.network.GetLocalIP()
	kademlia.network.Listen(ip, "8080")

	me := NewContact(id, ip)
	kademlia.RT = NewRoutingTable(me)
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	fmt.Println("hash = ", hash)
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	fmt.Println("hash = ", data)
	// TODO
}
