package gazz

type OctetString []byte

func (g OctetString) Encode(dst []byte) error {
	copy(dst, g)
	return nil
}

func (g OctetString) Len() int {
	return len(g)
}

func (g OctetString) Decode(src []byte) (any, error) {
	return src, nil
}
