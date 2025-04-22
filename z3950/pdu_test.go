package z3950

import (
	"encoding/asn1"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPDU(t *testing.T) {
	var pdu PDU

	ver := ProtocolVersion{
		BitLength: 8,
		Bytes:     []byte{0x7},
	}
	options := Options{
		BitLength: 8,
		Bytes:     []byte{0x7},
	}

	initRequest := InitRequest{
		ReferenceId:            nil,
		ProtocolVersion:        ver,
		Options:                options,
		PreferredMessageSize:   32768,
		ExceptionalMessageSize: 60000,
		ImplementationId:       []InternationalString{"81"},
		ImplementationName:     []InternationalString{"gazz"},
		ImplementationVersion:  []InternationalString{"0.0"},
	}
	pdu.InitRequest = append(pdu.InitRequest, initRequest)
	bytes, err := asn1.Marshal(pdu)
	assert.NoError(t, err, "Expected no error during ASN.1 marshaling")
	assert.NotEmpty(t, bytes, "Expected non-empty byte array after marshaling")

	//err = os.WriteFile("pdu.bin", bytes, 0644)
	//assert.NoError(t, err, "Expected no error during file write")

	for idx := range bytes {
		if idx%16 == 0 {
			fmt.Print("\n")
		}
		fmt.Printf("%02x ", bytes[idx])
	}
	var pdu2 PDU
	asn1.Unmarshal(bytes, &pdu2)
}
