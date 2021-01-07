package sexp

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
)

type Symbol string

type Cons struct {
	Car interface{}
	Cdr interface{}
}

func Marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	if err := e.Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

type Marshaler interface {
	MarshalSexp() ([]byte, error)
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
	case Marshaler:
		return enc.encodeMarshaler(v)
	case int, uint, int32, uint32, int64, uint64:
		return enc.printf("%d", v)
	case float32, float64:
		return enc.printf("%g", v)
	case Symbol:
		return enc.printf("%s", v)
	case string:
		return enc.printf("%q", v)
	case Cons:
		return enc.encodeCons(v)
	default:
		rv := reflect.ValueOf(v)
		// slices
		// arrays
		// structs? alist plist
		switch rv.Type().Kind() {
		case reflect.Ptr:
			return enc.Encode(rv.Elem().Interface())
		case reflect.Slice, reflect.Array:
			return enc.encodeList(rv)
		default:
			return fmt.Errorf("sexp encode: unsupported type %T", v)
		}

	}
}

func (enc *Encoder) encodeMarshaler(v Marshaler) error {
	b, err := v.MarshalSexp()
	if err != nil {
		return err
	}
	if _, err := enc.w.Write(b); err != nil {
		return err
	}
	return nil
}

func (enc *Encoder) encodeCons(v Cons) error {
	if err := enc.printf("("); err != nil {
		return err
	}
	if err := enc.Encode(v.Car); err != nil {
		return err
	}
	if err := enc.printf(" . "); err != nil {
		return err
	}
	if err := enc.Encode(v.Cdr); err != nil {
		return err
	}
	if err := enc.printf(")"); err != nil {
		return err
	}
	return nil
}

func (enc *Encoder) encodeList(v reflect.Value) error {
	l := v.Len()
	if l == 0 {
		return enc.printf("()")
	}
	if err := enc.printf("("); err != nil {
		return err
	}
	for i := 0; i < v.Len(); i++ {
		if i != 0 {
			if err := enc.printf(" "); err != nil {
				return err
			}
		}
		if err := enc.Encode(v.Index(i).Interface()); err != nil {
			return err
		}
	}
	if err := enc.printf(")"); err != nil {
		return err
	}
	return nil
}

func (enc *Encoder) printf(format string, v ...interface{}) error {
	_, err := fmt.Fprintf(enc.w, format, v...)
	return err
}
