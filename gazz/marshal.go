package gazz

import (
	"fmt"
	"reflect"
)

func EncodeTag(tclass int, tag int, constructed bool) ([]byte, error) {
	if tclass < 0 || tclass > 3 {
		return nil, fmt.Errorf("invalid class: %d", tclass)
	}
	var b byte = (byte(tclass) << 6)
	if constructed {
		b |= 0x20
	}
	if tag < 31 {
		if tag < 1 {
			return nil, fmt.Errorf("invalid tag: %d", tag)
		}
		b |= byte(tag)
		return []byte{b}, nil
	}
	l := 1
	tmpTag := tag
	for tmpTag > 0x7F {
		tmpTag >>= 7
		l++
	}
	dst := make([]byte, l+1)
	dst[0] = b
	for i := 1; i < l; i++ {
		b = byte(tag & 0x7F)
		if i != l-1 {
			b |= 0x00
		}
		dst[i] = b
		tag >>= 7
	}
	return dst, nil
}

func EncodeLen(len int) []byte {
	if len < 0 {
		return []byte{0x80}
	}
	if len < 0x80 {
		return []byte{byte(len)}
	}
	vl := Integer(len)
	n := vl.Len()
	dst := make([]byte, n+1)
	dst[0] = byte(0x80 | n)
	vl.Encode(dst[1:])
	return dst
}

type FieldParms struct {
	Tag      int
	Tclass   int
	Implicit bool
}

func Marshal(val any) ([]byte, error) {
	return MarshalTag(val, nil)
}

func MarshalTag(val any, fieldParms *FieldParms) ([]byte, error) {
	var dst []byte
	var codec Codec

	rv := reflect.ValueOf(val)
	if rv.Kind() == reflect.Struct {
		return nil, fmt.Errorf("struct types are not supported")
	}
	var tag int
	var tclass int
	switch v := val.(type) {
	case int:
		tag = TagInteger
		tclass = ClassUniversal
		codec = Integer(v)
	case []byte:
		tag = TagOctetString
		tclass = ClassUniversal
		codec = OctetString(v)
	default:
		return nil, fmt.Errorf("unsupported type: %T", val)
	}
	if fieldParms != nil {
		tag = fieldParms.Tag
		tclass = fieldParms.Tclass
	}
	dst, err := EncodeTag(tclass, tag, false)
	if err != nil {
		return nil, err
	}
	lenBytes := EncodeLen(codec.Len())
	dst = append(dst, lenBytes...)
	bodyBytes := make([]byte, codec.Len())
	err = codec.Encode(bodyBytes)
	if err != nil {
		return nil, err
	}
	dst = append(dst, bodyBytes...)
	return dst, err
}
