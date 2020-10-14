package kademlia

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSearchLocalStore(t *testing.T) {
	node := Node{nil, Client{}, make(map[string]string), 10}
	node.insertLocalStore("hello", "there")

	val1 := node.searchLocalStore("hello")
	assert.Equal(t, "there", *val1)

	val2 := node.searchLocalStore("shouldNotExist")
	assert.Nil(t, val2)

}

func TestUpdateContent(t *testing.T) {
	node := Node{nil, Client{}, make(map[string]string), 0}

	now := time.Now() // current local time
	sec := now.Unix() // number of seconds since January 1, 1970 UTC

	// create package and subtract 1000 seconds from current time to make it outdated
	data_package := strconv.FormatInt(sec-1000, 10) + ":" + "there"
	node.insertLocalStore("hello", data_package)
	val := node.content
	node.updateContent()

	assert.Equal(t, 0, len(val))

}

func TestGenerateRefreshNodeValue(t *testing.T) {
	assert.Equal(t, "0000000000000000000000000000000000000002", generateRefreshNodeValue(1, 1).String())
	assert.Equal(t, "0000000000000000000000000000000000000004", generateRefreshNodeValue(2, 1).String())
	assert.Equal(t, "0000000000000000000000000000010113070701", generateRefreshNodeValue(40, 1).String())
	assert.Equal(t, "00000000000000100b5e0038281912513b2f5751", generateRefreshNodeValue(100, 1).String())
	assert.Equal(t, "802935036b019b8104836f4026824e22449e125f", generateRefreshNodeValue(159, 1).String())
}
