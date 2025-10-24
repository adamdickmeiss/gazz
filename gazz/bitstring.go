package gazz

import "encoding/asn1"

type BitString asn1.BitString

func (g BitString) Encode(dst []byte) error {
	dst[0] = byte(g.BitLength % 8)
	for i := 0; i < len(g.Bytes); i++ {
		dst[i+1] = g.Bytes[i]
	}
	return nil
}

func (g BitString) Len() int {
	return len(g.Bytes) + 1
}

func (g BitString) Decode(src []byte) (any, error) {
	return src, nil
}
