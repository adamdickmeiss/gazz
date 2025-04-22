package gazz

type Codec interface {
	Encode([]byte) error
	Len() int
	Decode([]byte) (any, error)
}
