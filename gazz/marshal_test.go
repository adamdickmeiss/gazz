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

type MyString string
type StructNoTag2 struct {
	MyInt    int
	MyOctet  OctetString
	MyString MyString
}

func TestMarshalNoTag2(t *testing.T) {
	my := StructNoTag2{
		MyInt:    0x12345678,
		MyOctet:  OctetString{0x01, 0x02, 0x03},
		MyString: "bar",
	}
	checkMarshal(t, my)
}

type StructWithTag struct {
	MyInt    int         `asn1:"tag:5"`
	MyOctet  OctetString `asn1:"tag:31"`
	MyString string      `asn1:"tag:128,explicit"`
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

type StructWithOptionalAr struct {
	MyInt    []int  `asn1:"tag:5,optional"`
	MyString string `asn1:"tag:7"`
}

func TestMarshalWithOptionalAr(t *testing.T) {
	val := StructWithOptionalAr{MyString: "bar"}
	checkMarshal(t, val)

	val = StructWithOptionalAr{MyString: "bar", MyInt: []int{126}}
	checkMarshal(t, val)
}

type StructWithOptionalPtr struct {
	MyInt    *int   `asn1:"tag:5,optional"`
	MyString string `asn1:"tag:7"`
}

type StructWithOptional3 struct {
	MyInt    int    `asn1:"tag:5,optional"`
	MyString string `asn1:"tag:7"`
}

func TestMarshalWithOptionalPtr(t *testing.T) {
	val := StructWithOptionalPtr{MyString: "bar"}
	checkMarshal(t, val)

	myInt := 126
	val = StructWithOptionalPtr{MyString: "bar", MyInt: &myInt}
	vala := StructWithOptional3{MyString: "bar", MyInt: 126}
	checkMarshal2(t, vala, val)
}

type StructWithAsn1BitString struct {
	Options asn1.BitString `asn1:"tag:6"`
}

type StructWithJustBitString struct {
	Options BitString `asn1:"tag:6"`
}

func TestMarshalWithBitString(t *testing.T) {
	vala := StructWithAsn1BitString{Options: asn1.BitString{Bytes: []byte{0x01, 0x02}, BitLength: 12}}
	valb := StructWithJustBitString{Options: BitString{Bytes: []byte{0x01, 0x02}, BitLength: 12}}
	checkMarshal2(t, vala, valb)
	checkMarshal(t, vala)
}

type PDU struct {
	InitRequest *InitializeRequest `asn1:"tag:20,implicit"`
}

/*
InitializeRequest ::= SEQUENCE{
  referenceId                  ReferenceId OPTIONAL,
  protocolVersion              ProtocolVersion,
  options                      Options,
  preferredMessageSize   [5]   IMPLICIT INTEGER,
  exceptionalRecordSize  [6]   IMPLICIT INTEGER,
  idAuthentication       [7]   IdAuthentication OPTIONAL, -- see note below
  implementationId       [110] IMPLICIT InternationalString OPTIONAL,
  implementationName     [111] IMPLICIT InternationalString OPTIONAL,
  implementationVersion  [112] IMPLICIT InternationalString OPTIONAL,
  userInformationField   [11]  EXTERNAL OPTIONAL,
  otherInfo                    OtherInformation OPTIONAL
}
*/

type InitializeRequest struct {
	ReferenceId           ReferenceId    `asn1:"tag:2,optional"`
	ProtocolVersion       asn1.BitString `asn1:"tag:3"`
	Options               asn1.BitString `asn1:"tag:4"`
	PreferredMessageSize  int            `asn1:"tag:5"`
	ExceptionalRecordSize int            `asn1:"tag:6"`
	// IdAuthentication      IdAuthentication `asn1:"tag:7,optional"`
	Result                bool                `asn1:"tag:12"`
	ImplementationId      InternationalString `asn1:"tag:110,optional"`
	ImplementationName    InternationalString `asn1:"tag:112,optional"`
	ImplementationVersion InternationalString `asn1:"tag:112,optional"`
	// UserInformationField  *asn1.RawValue      `asn1:"tag:11,optional"`
	// OtherInformation      *OtherInformation   `asn1:"tag:200,optional"`
}

type IdAuthentication struct {
	// TODO
	Open string `asn1:"tag:2,explicit"`
}

type OtherInformation struct {
	Entries []OtherInformationEntry `asn1:"tag:201"`
}

type OtherInformationEntry struct {
	Category InfoCategory `asn1:"tag:202"`
}

type InfoCategory int

type InternationalString string // GeneralString
type ReferenceId OctetString    // [2]  IMPLICIT OCTET STRING

// not used as asn.1 package does not foll0ow custom BitString types
type ProtocolVersion asn1.BitString

// not used as asn.1 package does not foll0ow custom BitString types
type Options asn1.BitString

func TestMarshalZ3950InitalizeRequest(t *testing.T) {
	val := InitializeRequest{
		ReferenceId:           ReferenceId{0x30, 0x31, 0x32, 0x33},
		ProtocolVersion:       asn1.BitString{Bytes: []byte{0x37, 0x38}, BitLength: 16},
		Options:               asn1.BitString{Bytes: []byte{0xC0, 0x00}, BitLength: 16},
		PreferredMessageSize:  16384,
		ExceptionalRecordSize: 65536,
		Result:                true,
		ImplementationId:      "81",
		ImplementationName:    "Gazz-Go",
		ImplementationVersion: "1.0.0",
	}
	checkMarshal(t, val)
}
