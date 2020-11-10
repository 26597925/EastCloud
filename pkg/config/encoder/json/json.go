package json

import (
	"encoding/json"
	"github.com/26597925/EastCloud/pkg/config/encoder"
)

type jsonEncoder struct{}

func (j jsonEncoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (j jsonEncoder) Decode(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}

func (j jsonEncoder) Type() string {
	return "json"
}

func NewEncoder() encoder.Encoder {
	return jsonEncoder{}
}

