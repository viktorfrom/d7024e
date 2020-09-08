package network

import (
	"fmt"
	"log"
	"sync"
	"testing"

	"github.com/viktorfrom/d7024e-kademlia/internal/kademlia"
)

func TestPing(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		Listen("127.0.0.1", 9081)
		wg.Done()
	}()

	c := kademlia.NewContact(kademlia.NewKademliaID("FFFFFFFF00000000000000000000000000000000"), "127.0.0.1:9081")
	nt := Network{}

	fmt.Println("Sending PING...")
	resp, err := nt.SendPingMessage(&c)
	if err != nil {
		log.Fatal(err)
	}

	if resp == nil {
		log.Fatal("empty response")
	} else if *resp == "PONG" {
		fmt.Println("...PONG received!")
		wg.Done()
	}

	wg.Wait()
}
