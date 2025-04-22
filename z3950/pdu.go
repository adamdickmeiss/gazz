package z3950

import "encoding/asn1"

type PDU struct {
	InitRequest  []InitRequest  `asn1:"tag:20,implicit"`
	InitResponse []InitResponse `asn1:"tag:21,implicit"`
}

type InitRequest struct {
	ReferenceId            *ReferenceId `asn1:"optional"`
	ProtocolVersion        ProtocolVersion
	Options                Options
	PreferredMessageSize   int                   `asn1:"tag:5,implicit"`
	ExceptionalMessageSize int                   `asn1:"tag:6,implicit"`
	IdAuthentication       IdAuthentication      `asn1:"tag:7,optional"`
	ImplementationId       []InternationalString `asn1:"tag:110,implicit,optional"`
	ImplementationName     []InternationalString `asn1:"tag:111,implicit,optional"`
	ImplementationVersion  []InternationalString `asn1:"tag:112,implicit,optional"`
	OtherInfo              []InternationalString `asn1:"optional"`
}

type IdAuthentication struct {
	Open      []VisibleString
	IdPass    []IdPass
	Anonymous []asn1.RawValue
	Other     []External `asn1:"universal,tag:8,implicit"`
}

type IdPass struct {
	GroupId  []InternationalString `asn1:"tag:0,implicit,optional"`
	UserId   []InternationalString `asn1:"tag:1,implicit,optional"`
	Password []InternationalString `asn1:"tag:1,implicit,optional"`
}

/*
EXTERNAL  ::=  [UNIVERSAL 8] IMPLICIT SEQUENCE

	{
	direct-reference  OBJECT IDENTIFIER OPTIONAL,
	indirect-reference  INTEGER OPTIONAL,
	data-value-descriptor  ObjectDescriptor  OPTIONAL,
	encoding  CHOICE
	            {single-ASN1-type  [0] ANY,
	            octet-aligned     [1] IMPLICIT OCTET STRING,
	            arbitrary         [2] IMPLICIT BIT STRING}
	}
*/
type External struct {
	DirectReference   []asn1.ObjectIdentifier `asn1:"optional"`
	IndirectReference []int                   `asn1:"optional"`
	ObjectDescriptor  ObjectDescriptor        `asn1:"optional"`
	Encoding          ExternalEncoding
}

type ExternalEncoding struct {
	SingleASN1Type []asn1.RawValue  `asn1:"tag:0"`
	OctetAligned   []OctetString    `asn1:"tag:1,implicit"`
	Arbitrary      []asn1.BitString `asn1:"tag:2,implicit"`
}

/*
-- OtherInformation

	OtherInformation   ::= [201] IMPLICIT SEQUENCE OF SEQUENCE{
	  category            [1]   IMPLICIT InfoCategory OPTIONAL,
	  information        CHOICE{
	    characterInfo        [2]  IMPLICIT InternationalString,
	    binaryInfo        [3]  IMPLICIT OCTET STRING,
	    externallyDefinedInfo    [4]  IMPLICIT EXTERNAL,
	    oid          [5]  IMPLICIT OBJECT IDENTIFIER}}
*/
type OtherInformation string // TODO: Define the structure
type OctetString []byte
type ObjectDescriptor string
type InternationalString string
type VisibleString string
type ProtocolVersion asn1.BitString
type Options asn1.BitString

type InitResponse struct {
	ReferenceId *ReferenceId `asn1:"optional"`
}

type ReferenceId string

func NewPDU() *PDU {
	return &PDU{}
}
