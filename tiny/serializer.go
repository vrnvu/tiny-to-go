package tiny

type Serializer interface {
	Decode(input []byte) (*Redirect, error)
	Encode(input *Redirect) ([]byte, error)
}
