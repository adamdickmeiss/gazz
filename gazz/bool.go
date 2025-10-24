package gazz

import "fmt"

type Bool bool

func (b Bool) Encode(dst []byte) error {
	if b {
		dst[0] = 0xFF
	} else {
		dst[0] = 0x00
	}
	return nil
}

func (b Bool) Len() int {
	return 1
}

func (b Bool) Decode(src []byte) (any, error) {
	if len(src) == 0 {
		return false, fmt.Errorf("empty boolean")
	}
	return src[0] != 0, nil
}
