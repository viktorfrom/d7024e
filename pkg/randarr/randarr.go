package randarr

import (
	"math/rand"
	"time"
)

// RandomBytes returns a byte slice with n number of random bytes
func RandomBytes(n int) []byte {
	rand.Seed(time.Now().UnixNano())
	arr := make([]byte, n)
	rand.Read(arr)
	return arr
}
