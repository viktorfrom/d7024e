package main

import (
	"strconv"
	"testing"
)

func TestAdd(t *testing.T) {
	res := add(10, 5)
	if res != 15 {
		t.Error("Expected 10 + 5 to equal 15 got " + strconv.FormatInt(int64(res), 10))
	}
}
