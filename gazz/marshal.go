package gazz

import (
	"fmt"
	"reflect"
)

type MarshalError struct {
	Message string
}

func (e MarshalError) Error() string {
	return e.Message
}

func EncodeTag(tclass int, tag int, constructed bool) ([]byte, error) {
	if tclass < 0 || tclass > 3 {
		return nil, &MarshalError{Message: fmt.Sprintf("invalid class: %d", tclass)}
	}
	var b byte = (byte(tclass) << 6)
	if constructed {
		b |= 0x20
	}
	if tag < 31 {
		if tag < 1 {
			return nil, &MarshalError{Message: fmt.Sprintf("invalid tag: %d", tag)}
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

func Marshal(val any) ([]byte, error) {
	return MarshalTag(val, nil)
}

func MarshalSequence(rv reflect.Value) ([]byte, error) {
	var bodyBytes []byte
	for i := 0; i < rv.NumField(); i++ {
		field := rv.Type().Field(i)
		fp, err := ParseStructField(field)
		if err != nil {
			return nil, err
		}
		if rv.Field(i).Kind() == reflect.Ptr || rv.Field(i).Kind() == reflect.Slice {
			if rv.Field(i).IsNil() {
				if fp == nil || !fp.Optional {
					return nil, &MarshalError{Message: fmt.Sprintf("field %s is nil and not optional", field.Name)}
				}
				continue
			}
		}
		fieldValue := rv.Field(i)
		encodedField, err := MarshalTag(fieldValue.Interface(), fp)
		if err != nil {
			return nil, err
		}
		bodyBytes = append(bodyBytes, encodedField...)
	}
	dst, err := EncodeTag(ClassUniversal, TagSequence, true)
	if err != nil {
		return nil, err
	}
	lenBytes := EncodeLen(len(bodyBytes))
	dst = append(dst, lenBytes...)
	dst = append(dst, bodyBytes...)
	return dst, nil
}

func MarshalChoice(rv reflect.Value) ([]byte, error) {
	for i := 0; i < rv.NumField(); i++ {
		fieldValue := rv.Field(i)
		if fieldValue.Kind() == reflect.Ptr && !fieldValue.IsNil() {
			field := rv.Type().Field(i)
			fp, err := ParseStructField(field)
			if err != nil {
				return nil, err
			}
			return MarshalTag(fieldValue.Elem().Interface(), fp)
		}
	}
	return nil, &MarshalError{Message: "no valid choice found"}
}

func MarshalTag(val any, fieldParms *StructTags) ([]byte, error) {
	if fieldParms != nil && fieldParms.Explicit {
		dst1, err := MarshalTag(val, nil)
		if err != nil {
			return nil, err
		}
		dst, err := EncodeTag(fieldParms.Tclass, fieldParms.Tag, true)
		if err != nil {
			return nil, err
		}
		lenBytes := EncodeLen(len(dst1))
		dst = append(dst, lenBytes...)
		dst = append(dst, dst1...)
		return dst, nil
	}
	var codec Codec
	var tag int
	var tclass int
	rv := reflect.ValueOf(val)
	switch rv.Kind() {
	case reflect.Struct:
		allPointers := true
		for i := 0; i < rv.NumField(); i++ {
			fieldValue := rv.Field(i)
			if fieldValue.Kind() != reflect.Ptr {
				allPointers = false
			}
		}
		if allPointers {
			return MarshalChoice(rv)
		} else {
			return MarshalSequence(rv)
		}
	case reflect.Slice:
		if v, ok := val.([]byte); ok {
			tag = TagOctetString
			tclass = ClassUniversal
			codec = OctetString(v)
		} else if v, ok := val.(OctetString); ok {
			tag = TagOctetString
			tclass = ClassUniversal
			codec = v
		} else {
			if rv.Len() == 0 {
				return nil, nil
			}
			if rv.Len() == 1 {
				fieldParms.Explicit = true
				return MarshalTag(rv.Index(0).Interface(), fieldParms)
			}
			return nil, &MarshalError{Message: fmt.Sprintf("unsupported type: %T", val)}
		}
	case reflect.String:
		tag = TagPrintableString
		tclass = ClassUniversal
		codec = OctetString(val.(string))
	case reflect.Int:
		tag = TagInteger
		tclass = ClassUniversal
		codec = Integer(val.(int))
	default:
		return nil, &MarshalError{Message: fmt.Sprintf("unsupported type: %T", val)}
	}
	if fieldParms != nil {
		tclass = fieldParms.Tclass
		tag = fieldParms.Tag
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
