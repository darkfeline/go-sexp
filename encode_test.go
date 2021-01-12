// Copyright (C) 2021 Allen Li
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
	t.Run("list", func(t *testing.T) {
		t.Parallel()
		run(t, testcase{v: []interface{}{5, "shiori"}, want: `(5 "shiori")`})
	})
	t.Run("alist", func(t *testing.T) {
		t.Parallel()
		t.Run("without coding", func(t *testing.T) {
			t.Parallel()
			type d struct {
				Princess string
			}
			run(t, testcase{v: d{"yui"}, want: `((Princess . "yui"))`})
		})
		t.Run("with coding", func(t *testing.T) {
			t.Parallel()
			type d struct {
				_sexpCoding struct{} `alist`
				Pri         string
			}
			run(t, testcase{v: d{Pri: "yui"}, want: `((Pri . "yui"))`})
		})
		t.Run("named field", func(t *testing.T) {
			t.Parallel()
			type d struct {
				_sexpCoding struct{} `alist`
				Pri         string   `sexp:"princess"`
			}
			run(t, testcase{v: d{Pri: "yui"}, want: `((princess . "yui"))`})
		})
	})
	t.Run("plist", func(t *testing.T) {
		t.Parallel()
		t.Run("default", func(t *testing.T) {
			t.Parallel()
			type d struct {
				_sexpCoding struct{} `plist`
				Pri         string
			}
			run(t, testcase{v: d{Pri: "yui"}, want: `(Pri "yui")`})
		})
		t.Run("named field", func(t *testing.T) {
			t.Parallel()
			type d struct {
				_sexpCoding struct{} `plist`
				Pri         string   `sexp:"princess"`
			}
			run(t, testcase{v: d{Pri: "yui"}, want: `(princess "yui")`})
		})
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
