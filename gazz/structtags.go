package gazz

import (
	"reflect"
	"strconv"
	"strings"
)

type StructTags struct {
	Tag      int
	Tclass   int
	Explicit bool
	Optional bool
}

func ParseStructField(f reflect.StructField) (*StructTags, error) {
	if f.Anonymous {
		return nil, nil
	}
	return ParseStructTag(f.Tag.Get("asn1"))
}

func ParseStructTag(s string) (*StructTags, error) {
	if s == "" {
		return nil, nil
	}
	var tags StructTags
	tags.Tag = 0
	tags.Tclass = ClassContextSpecific
	tags.Explicit = false
	tags.Optional = false
	for _, comp := range strings.Split(s, ",") {
		switch {
		case strings.HasPrefix(comp, "tag:"):
			no, err := strconv.Atoi(strings.TrimPrefix(comp, "tag:"))
			if err != nil {
				return nil, err
			}
			tags.Tag = no
		case comp == "explicit":
			tags.Explicit = true
		case comp == "optional":
			tags.Optional = true
		case comp == "application":
			tags.Tclass = ClassApplication
		case comp == "private":
			tags.Tclass = ClassPrivate
		}
	}
	return &tags, nil
}
