package sexp

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strings"
)

// A Symbol is marshaled as a symbol.
type Symbol string

// A Cons is marshaled as a cons pair.
type Cons struct {
	Car interface{}
	Cdr interface{}
}

// Marshal marshals a value into an s expression.
func Marshal(v interface{}) ([]byte, error) {
	var b bytes.Buffer
	e := NewEncoder(&b)
	if err := e.Encode(v); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

// A Marshaler defines how to marshal itself as an s expression.
type Marshaler interface {
	MarshalSexp() ([]byte, error)
}

// An Encoder encodes values as s expressions.
type Encoder struct {
	w   io.Writer
	err error
}

// NewEncoder creates a new Encoder.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: w,
	}
}

// Encode a value as an s expression.
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
		// structs? alist plist
		switch rv.Type().Kind() {
		case reflect.Ptr:
			return enc.Encode(rv.Elem().Interface())
		case reflect.Slice, reflect.Array:
			return enc.encodeList(rv)
		case reflect.Struct:
			return enc.encodeStruct(rv)
		default:
			enc.err = fmt.Errorf("sexp encode: unsupported type %T", v)
			return enc.err
		}
	}
}

func (enc *Encoder) encodeMarshaler(v Marshaler) error {
	if enc.err != nil {
		return enc.err
	}
	b, err := v.MarshalSexp()
	if err != nil {
		enc.err = err
		return err
	}
	if _, err := enc.w.Write(b); err != nil {
		enc.err = err
		return err
	}
	return nil
}

func (enc *Encoder) encodeCons(v Cons) error {
	enc.printf("(")
	enc.Encode(v.Car)
	enc.printf(" . ")
	enc.Encode(v.Cdr)
	enc.printf(")")
	return enc.err
}

func (enc *Encoder) encodeList(v reflect.Value) error {
	l := v.Len()
	if l == 0 {
		return enc.printf("()")
	}
	enc.printf("(")
	for i := 0; i < v.Len(); i++ {
		if i != 0 {
			enc.printf(" ")
		}
		enc.Encode(v.Index(i).Interface())
	}
	enc.printf(")")
	return enc.err
}

func (enc *Encoder) encodeStruct(v reflect.Value) error {
	f, ok := v.Type().FieldByName("_sexpCoding")
	if !ok || f.Tag == "alist" {
		return enc.encodeStructAlist(v)
	}
	if f.Tag == "plist" {
		return enc.encodeStructPlist(v)
	}
	enc.err = fmt.Errorf("sexp encode: struct %T with bad _sexpCoding tag %s",
		v.Interface(), f.Tag)
	return enc.err
}

func (enc *Encoder) encodeStructAlist(v reflect.Value) error {
	enc.printf("(")
	first := true
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			// Unexported field
			continue
		}
		if first {
			first = false
		} else {
			enc.printf(" ")
		}
		enc.Encode(Cons{fieldKey(f), v.Field(i).Interface()})
	}
	enc.printf(")")
	return enc.err
}

func (enc *Encoder) encodeStructPlist(v reflect.Value) error {
	enc.printf("(")
	first := true
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.PkgPath != "" {
			// Unexported field
			continue
		}
		if first {
			first = false
		} else {
			enc.printf(" ")
		}
		enc.Encode(fieldKey(f))
		enc.printf(" ")
		enc.Encode(v.Field(i).Interface())
	}
	enc.printf(")")
	return enc.err
}

func (enc *Encoder) printf(format string, v ...interface{}) error {
	if enc.err != nil {
		return enc.err
	}
	_, err := fmt.Fprintf(enc.w, format, v...)
	enc.err = err
	return err
}

// Return key to use for field
func fieldKey(f reflect.StructField) Symbol {
	tag := f.Tag.Get("sexp")
	parts := strings.SplitN(tag, ",", 2)
	if name := parts[0]; name != "" {
		return Symbol(name)
	}
	return Symbol(f.Name)
}
