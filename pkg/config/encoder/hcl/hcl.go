package hcl

import (
	"encoding/json"
	"github.com/hashicorp/hcl"
	"sapi/pkg/config/encoder"
)

type hclEncoder struct{}

func (h hclEncoder) Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func (h hclEncoder) Decode(d []byte, v interface{}) error {
	return hcl.Unmarshal(d, v)
}

func (h hclEncoder) Type() string {
	return "hcl"
}

func NewEncoder() encoder.Encoder {
	return hclEncoder{}
}

