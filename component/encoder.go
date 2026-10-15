package component

import (
	"encoding/json"
	"io"
)

// Marshal encodes a component using the standard library JSON encoder.
func Marshal(c Component) ([]byte, error) {
	return json.Marshal(c)
}

// Encoder writes components to a stream using the standard library encoder.
type Encoder struct {
	encoder *json.Encoder
}

// NewEncoder creates a streaming component encoder.
func NewEncoder(writer io.Writer) *Encoder {
	return &Encoder{encoder: json.NewEncoder(writer)}
}

// Encode writes one component followed by a newline.
func (e *Encoder) Encode(c Component) error {
	return e.encoder.Encode(c)
}
