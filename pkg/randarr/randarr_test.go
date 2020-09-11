package randarr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCorrectLength(t *testing.T) {
	zero := RandomBytes(0)
	ten := RandomBytes(10)
	twenty := RandomBytes(20)
	thirty := RandomBytes(30)

	assert.Equal(t, 0, len(zero))
	assert.Equal(t, 10, len(ten))
	assert.Equal(t, 20, len(twenty))
	assert.Equal(t, 30, len(thirty))
}
