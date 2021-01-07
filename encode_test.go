package sexp

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestEncode(t *testing.T) {
	t.Parallel()
	t.Run("literals", func(t *testing.T) {
		t.Parallel()
		t.Run("string", func(t *testing.T) {
			t.Parallel()
			got, err := Marshal("ionasal")
			if err != nil {
				t.Fatal(err)
			}
			const want = `"ionasal"`
			if diff := cmp.Diff(want, string(got)); diff != "" {
				t.Errorf("output mismatch (-want +got):\n%s", diff)
			}
		})
	})
}
