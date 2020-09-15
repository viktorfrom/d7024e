package randarr

import (
	"encoding/hex"
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

// RandomHexString returns a random hexadecimal string with the size n bytes.
func RandomHexString(n int) string {
	bytes := RandomBytes(n)
	encodedString := hex.EncodeToString(bytes)
	return encodedString
}
