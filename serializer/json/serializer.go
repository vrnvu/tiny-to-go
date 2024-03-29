package json

import (
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/vrnvu/tiny-to-go/tiny"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*tiny.Redirect, error) {
	redirect := &tiny.Redirect{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return nil, errors.Wrap(err, "serialier.Redirect.Decode")
	}
	return redirect, nil
}

func (r *Redirect) Encode(input *tiny.Redirect) ([]byte, error) {
	raw, err := json.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serialier.Redirect.Encode")
	}
	return raw, nil
}
