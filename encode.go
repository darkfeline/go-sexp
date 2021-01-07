package sexp

import (
	"bytes"
	"io"
)

type Symbol string

func Marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	if err := e.Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

type Encoder struct{}

func NewEncoder(w io.Writer) *Encoder {}

func (enc *Encoder) Encode(v interface{}) error {}
