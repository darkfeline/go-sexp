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
		t.Run("string", func(t *testing.T) {
			t.Parallel()
			run(t, testcase{v: "ionasal", want: `"ionasal"`})
		})
		t.Run("escaped string", func(t *testing.T) {
			t.Parallel()
			run(t, testcase{v: `"ionasal"`, want: `"\"ionasal\""`})
		})
	})
}
