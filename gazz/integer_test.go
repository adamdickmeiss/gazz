package gazz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIntegerEncode12345678(t *testing.T) {
	var i Integer = 0x12345678
	assert.Equal(t, 4, i.Len())
	dst := make([]byte, i.Len())
	err := i.Encode(dst)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x12, 0x34, 0x56, 0x78}, dst)
}

func TestIntegerEncodeFFFF(t *testing.T) {
	var i Integer = 0xFFFF
	assert.Equal(t, 2, i.Len())
	dst := make([]byte, i.Len())
	err := i.Encode(dst)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0xFF, 0xFF}, dst)
}

func TestIntegerEncode0(t *testing.T) {
	var i Integer = 0x0
	assert.Equal(t, 1, i.Len())
	dst := make([]byte, i.Len())
	err := i.Encode(dst)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x0}, dst)
}
