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
	case Symbol:
		_, err := fmt.Fprintf(enc.w, "%s", v)
		return err
	case string:
		_, err := fmt.Fprintf(enc.w, "%q", v)
		return err
	case int, uint, int32, uint32, int64, uint64:
		_, err := fmt.Fprintf(enc.w, "%d", v)
		return err
	case float32, float64:
		_, err := fmt.Fprintf(enc.w, "%g", v)
		return err
	default:
		return fmt.Errorf("sexp encode: unsupported type %T", v)
	}
}
