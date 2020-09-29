package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const KADEMLIA_BIN string = "./kademlia"

func Test(t *testing.T) {
	os.Args = []string{"./calc", "-mode", "multiply", "3", "2", "5"}
	assert.Equal(t, 1, 1)
}
