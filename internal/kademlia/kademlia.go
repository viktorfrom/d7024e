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

func (kademlia *Kademlia) Ping(input string) {
	target := NewContact(NewRandomKademliaID(), input+":8080")
	rpc, err := kademlia.network.SendPingMessage(&target, &kademlia.RT.me)

	if err != nil {
		fmt.Println(err)
	} else if *rpc.Type == "OK" {
		kademlia.RT.AddContact(target)
		fmt.Println("Add " + target.ID.String() + " to bucket id: " + strconv.Itoa(kademlia.RT.getBucketIndex(target.ID)))
	} else {
		fmt.Print("CAOS?")
	}
}
