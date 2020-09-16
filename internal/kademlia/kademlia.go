package kademlia

import (
	"fmt"
	"strconv"
)

type Kademlia struct {
	RT      *RoutingTable
	network Network
}

// InitNode initializes the Kademlia Node
// with a Routing Table and a Network
func (kademlia *Kademlia) InitNode(id *KademliaID) {
	kademlia.network = Network{kademlia}
	ip := kademlia.network.GetLocalIP()
	go kademlia.network.Listen(ip, "8080")

	me := NewContact(id, ip+":8080")
	rendezvousNode := NewContact(NewKademliaID("00000000000000000000000000000000FFFFFFFF"), "10.0.8.3:8080")
	kademlia.RT = NewRoutingTable(me)
	kademlia.RT.AddContact(rendezvousNode)
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

func (kademlia *Kademlia) Ping() {
	target := &kademlia.RT.FindClosestContacts(kademlia.RT.me.ID, bucketSize)[0]
	rpc, err := kademlia.network.SendPingMessage(target, &kademlia.RT.me)

	fmt.Println("Send ping")

	if err != nil {
		fmt.Println(err)
		kademlia.RT.RemoveContact(*target)
	} else if *rpc.Type == "OK" {
		fmt.Println("Add " + target.ID.String() + " to bucket id: " + strconv.Itoa(kademlia.RT.getBucketIndex(target.ID)))
		kademlia.RT.AddContact(*target)
	} else {
		fmt.Print("CAOS?")
	}
}
