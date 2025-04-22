package gazz

import (
	"encoding/asn1"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalInteger(t *testing.T) {
	var i int = 0x12345678
	bytes, err := Marshal(i)
	assert.NoError(t, err)
	var expected []byte = []byte{0x02, 0x04, 0x12, 0x34, 0x56, 0x78}
	assert.Equal(t, expected, bytes)

	bytes, err = asn1.Marshal(i)
	assert.NoError(t, err)
	assert.Equal(t, expected, bytes)
}

type StructNoTag struct {
	MyInt    int
	MyOctet  OctetString
	MyString string
}

func TestMarshalNoTag(t *testing.T) {
	my := StructNoTag{
		MyInt:    0x12345678,
		MyOctet:  OctetString{0x01, 0x02, 0x03},
		MyString: "bar",
	}
	bytes, err := asn1.Marshal(my)
	assert.NoError(t, err)
	assert.NotEmpty(t, bytes)
	assert.Equal(t, []byte{0x30, 0x10, 0x2, 0x4, 0x12, 0x34, 0x56, 0x78, 0x4, 0x3, 0x1, 0x2, 0x3, 0x13, 0x03, 0x62, 0x61, 0x72}, bytes)
}

type StructWithTag struct {
	MyInt    int         `asn1:"tag:5"`
	MyOctet  OctetString `asn1:"tag:6"`
	MyString string      `asn1:"tag:7,implicit"`
}

func TestMarshalWithTag(t *testing.T) {
	my := StructWithTag{
		MyInt:    0x12345678,
		MyOctet:  OctetString{0x01, 0x02, 0x03},
		MyString: "bar",
	}
	bytes, err := asn1.Marshal(my)
	assert.NoError(t, err)
	assert.NotEmpty(t, bytes)
	assert.Equal(t, []byte{0x30, 0x10, 0x85, 0x4, 0x12, 0x34, 0x56, 0x78, 0x86, 0x3, 0x1, 0x2, 0x3, 0x87, 0x03, 0x62, 0x61, 0x72}, bytes)
}
