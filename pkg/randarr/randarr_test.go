package randarr

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteArrayCorrectLength(t *testing.T) {
	zero := RandomBytes(0)
	ten := RandomBytes(10)
	twenty := RandomBytes(20)
	thirty := RandomBytes(30)

	assert.Equal(t, 0, len(zero))
	assert.Equal(t, 10, len(ten))
	assert.Equal(t, 20, len(twenty))
	assert.Equal(t, 30, len(thirty))
}

func TestHexStrCorrectLength(t *testing.T) {
	zeroHex := RandomHexString(0)
	tenHex := RandomHexString(10)
	twentyHex := RandomHexString(20)
	thirtyHex := RandomHexString(30)

	ten, _ := hex.DecodeString(tenHex)
	twenty, _ := hex.DecodeString(twentyHex)
	thirty, _ := hex.DecodeString(thirtyHex)

	assert.Equal(t, 0, len(zeroHex))
	assert.Equal(t, 10, len(ten))
	assert.Equal(t, 20, len(twenty))
	assert.Equal(t, 30, len(thirty))
}
