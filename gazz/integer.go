package gazz

import "fmt"

type Integer int

func (i Integer) Encode(dst []byte) error {
	n := i.Len()
	for j := 0; n > 0; j++ {
		n--
		dst[j] = byte(i >> uint(n*8))
	}
	return nil
}

func (i Integer) Len() int {
	len := 1
	for i > 255 {
		i >>= 8
		len++
	}
	return len
}

func (i Integer) Decode(src []byte) (any, error) {
	if len(src) == 0 {
		return 0, fmt.Errorf("empty integer")
	}
	var n int
	for _, b := range src {
		n = (n << 8) | int(b)
	}
	return n, nil
}
