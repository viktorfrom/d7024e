package kademlia

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchLocalStore(t *testing.T) {
	node := Node{nil, Client{}, make(map[string]string), 10}
	node.insertLocalStore("hello", "there")

	c := node.searchLocalStore("hello")
	assert.Equal(t, "there", *c)

	y := node.searchLocalStore("shouldNotExist")
	assert.Nil(t, y)

}

func TestGenerateRefreshNodeValue(t *testing.T) {
	assert.Equal(t, "0000000000000000000000000000000000000002", generateRefreshNodeValue(1, 1).String())
	assert.Equal(t, "0000000000000000000000000000000000000004", generateRefreshNodeValue(2, 1).String())
	assert.Equal(t, "0000000000000000000000000000010113070701", generateRefreshNodeValue(40, 1).String())
	assert.Equal(t, "00000000000000100b5e0038281912513b2f5751", generateRefreshNodeValue(100, 1).String())
	assert.Equal(t, "802935036b019b8104836f4026824e22449e125f", generateRefreshNodeValue(159, 1).String())
}
