package sexp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEncode(t *testing.T) {
	t.Parallel()
	type testcase struct {
		v    interface{}
		want string
	}
	run := func(t *testing.T, tc testcase) {
		got, err := Marshal(tc.v)
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(tc.want, string(got)); diff != "" {
			t.Errorf("output mismatch (-want +got):\n%s", diff)
		}
	}
	t.Run("literals", func(t *testing.T) {
		t.Parallel()
		t.Run("int", func(t *testing.T) {
			t.Parallel()
			run(t, testcase{v: 5, want: `5`})
		})
		t.Run("float", func(t *testing.T) {
			t.Parallel()
			run(t, testcase{v: 4.25, want: `4.25`})
		})
		t.Run("string", func(t *testing.T) {
			t.Parallel()
			run(t, testcase{v: "ionasal", want: `"ionasal"`})
		})
		t.Run("escaped string", func(t *testing.T) {
			t.Parallel()
			run(t, testcase{v: `"ionasal"`, want: `"\"ionasal\""`})
		})
		t.Run("symbol", func(t *testing.T) {
			t.Parallel()
			run(t, testcase{v: Symbol("1+"), want: `1+`})
		})
	})
	t.Run("cons", func(t *testing.T) {
		t.Parallel()
		run(t, testcase{v: Cons{1, 2}, want: `(1 . 2)`})
	})
	t.Run("pointer", func(t *testing.T) {
		t.Parallel()
		v := 5
		run(t, testcase{v: &v, want: `5`})
	})
	t.Run("marshaler", func(t *testing.T) {
		t.Parallel()
		run(t, testcase{
			v:    testMarshaler{[]byte(`(kokkoro peco kyaru)`)},
			want: `(kokkoro peco kyaru)`,
		})
	})
}

type testMarshaler struct {
	b []byte
}

func (m testMarshaler) MarshalSexp() ([]byte, error) {
	return m.b, nil
}
