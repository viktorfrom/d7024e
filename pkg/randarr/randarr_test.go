package randarr

import (
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
	zero, _ := RandomHexString(0)
	ten, _ := RandomHexString(10)
	twenty, _ := RandomHexString(20)
	thirty, _ := RandomHexString(30)

	assert.Equal(t, 0, len(zero))
	assert.Equal(t, 10, len([]byte(ten)))
	assert.Equal(t, 20, len([]byte(twenty)))
	assert.Equal(t, 30, len([]byte(thirty)))
}

func TestHexStrOddNumber(t *testing.T) {
	_, err := RandomHexString(3)
	assert.Error(t, err)
}
