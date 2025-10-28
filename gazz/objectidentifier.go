package gazz

import "encoding/asn1"

type ObjectIdentifier asn1.ObjectIdentifier

func (o ObjectIdentifier) Encode(dst []byte) error {
	n := o.Len()
	for j := 0; n > 0; j++ {
		n--
		dst[j] = byte(o[j])
	}
	return nil
}

func (o ObjectIdentifier) Len() int {
	return len(o)
}

func (o ObjectIdentifier) Decode(src []byte) (any, error) {
	var oid asn1.ObjectIdentifier
	for _, b := range src {
		oid = append(oid, int(b))
	}
	return oid, nil
}
