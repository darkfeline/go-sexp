package sexp

import (
	"bytes"
	"fmt"
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

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

func (enc *Encoder) Encode(v interface{}) error {
	switch v := v.(type) {
	case string:
		fmt.Fprintf(enc.w, "%q", v)
		return nil
	case int, uint, int32, uint32, int64, uint64:
		fmt.Fprintf(enc.w, "%d", v)
		return nil
	default:
		return fmt.Errorf("sexp encode: unsupported type %T", v)
	}
}
