package randarr

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

const (
	errUneven string = "n is not an even number"
)

// RandomBytes returns a byte slice with n number of random bytes
func RandomBytes(n int) []byte {
	rand.Seed(time.Now().UnixNano())
	arr := make([]byte, n)
	rand.Read(arr)
	return arr
}

// RandomHexString returns a hexadecimal string with the size n bytes.
// n has to be an even number otherwise an error will be returned
func RandomHexString(n int) (string, error) {
	if n%2 != 0 {
		return "", errors.New(errUneven)
	}
	n = n / 2
	arr := RandomBytes(n)
	fmt.Println("arr: ", arr)
	encodedString := hex.EncodeToString(arr)
	return encodedString, nil
}
