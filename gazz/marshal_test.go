package gazz

import (
	"encoding/asn1"
	"testing"

	"github.com/stretchr/testify/assert"
)

func checkMarshal2(t *testing.T, vala any, valb any) {
	expected, err := asn1.Marshal(vala)
	assert.NoError(t, err)

	bytes, err := Marshal(valb)
	assert.NoError(t, err)

	assert.Equal(t, expected, bytes)
}

func checkMarshal(t *testing.T, val any) {
	checkMarshal2(t, val, val)
}

func TestMarshalInteger(t *testing.T) {
	var i int = 0x12345678
	checkMarshal(t, i)
}

func TestMarshalString(t *testing.T) {
	var s string = "hello"
	checkMarshal(t, s)
}

func TestMarshalByteSlice(t *testing.T) {
	var b []byte = []byte{0x97, 0x98}
	checkMarshal(t, b)
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
	checkMarshal(t, my)
}

type StructWithTag struct {
	MyInt    int         `asn1:"tag:5"`
	MyOctet  OctetString `asn1:"tag:6"`
	MyString string      `asn1:"tag:7,explicit"`
}

func TestMarshalWithTag(t *testing.T) {
	my := StructWithTag{
		MyInt:    0x12345678,
		MyOctet:  OctetString{0x01, 0x02, 0x03},
		MyString: "bar",
	}
	checkMarshal(t, my)
}

type StructWithChoice struct {
	MyInt    *int    `asn1:"tag:5"`
	MyString *string `asn1:"tag:7"`
}

func TestMarshalWithChoice1(t *testing.T) {
	myint := 3
	val := StructWithChoice{MyInt: &myint}

	bytes, err := Marshal(val)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x85, 0x1, 0x3}, bytes)
}

func TestMarshalWithChoice2(t *testing.T) {
	bar := "bar"
	val := StructWithChoice{MyString: &bar}

	bytes, err := Marshal(val)
	assert.NoError(t, err)
	assert.Equal(t, []byte{0x87, 0x03, 0x62, 0x61, 0x72}, bytes)
}

type StructWithOptional1 struct {
	MyInt    []int  `asn1:"tag:5,optional"`
	MyString string `asn1:"tag:7"`
}

func TestMarshalWithOptional1(t *testing.T) {
	val := StructWithOptional1{MyString: "bar"}
	checkMarshal(t, val)

	val = StructWithOptional1{MyString: "bar", MyInt: []int{126}}
	checkMarshal(t, val)
}

type StructWithBitStringA struct {
	Options asn1.BitString `asn1:"tag:6"`
}

type StructWithBitStringB struct {
	Options BitString `asn1:"tag:6"`
}

func TestMarshalWithBitString(t *testing.T) {
	vala := StructWithBitStringA{Options: asn1.BitString{Bytes: []byte{0x01, 0x02}, BitLength: 12}}
	valb := StructWithBitStringB{Options: BitString{Bytes: []byte{0x01, 0x02}, BitLength: 12}}
	checkMarshal2(t, vala, valb)
}
