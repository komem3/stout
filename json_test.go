package main

import (
	"bytes"
	"encoding/json"
	"go/ast"
	"stout/internal/testutil"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStType2Json(t *testing.T) {
	type (
		given struct {
			path   string
			stType string
		}
		want struct {
			outStruct interface{}
			err       error
		}
	)
	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			"SampleJson",
			given{
				path:   "./define_test.go",
				stType: "SampleJson",
			},
			want{
				outStruct: SampleJson{
					ID:      1,
					Content: "string",
					Tag:     1,
					Ignore:  1,
					Hyphen:  1,
					PStr:    pStr("string"),
					Pointer: &Embbading{
						IDString: "string",
					},
					Struct: Embbading{
						IDString: "string",
					},
					Array: []string{"string"},
					ArraySct: []Embbading{
						{
							IDString: "string",
						},
					},
					ArrayPSct: []*Embbading{
						{
							IDString: "string",
						},
					},
					Other: ast.ArrayType{
						Lbrack: 1,
						Len:    nil,
						Elt:    nil,
					},
					Internal: testutil.OtherPkg{
						PkgContent: "string",
					},
					Uinter: 1,
					Embbading: Embbading{
						IDString: "string",
					},
					Integer: 1,
					ArrayEmbbad: ArrayEmbbad{
						{
							IDString: "string",
						},
					},
					PString: func(str PString) *PString { return &str }("string"),
					PEmb: &PEmb{
						Content: "string",
						Next:    "string",
					},
					OtherPkg: &testutil.OtherPkg{
						PkgContent: "string",
					},
					ArrayPEmbbad: &ArrayPEmbbad{
						{
							PkgContent: "string",
						},
					},
					private: "private",
				},
				err: nil,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var testFunc = func(t *testing.T, format bool) {
				var (
					got  = new(bytes.Buffer)
					want = new(bytes.Buffer)
				)
				encoder := json.NewEncoder(want)
				if format {
					encoder.SetIndent("", "  ")
				}

				if err := encoder.Encode(tt.want.outStruct); err != nil {
					t.Fatal(err)
				}

				op := newJsonOption(tt.given.path, tt.given.stType, !format)
				err := stType2Json(got, op)
				if diff := cmp.Diff(err, tt.want.err, cmpopts.EquateErrors()); diff != "" {
					t.Errorf("error: got(-) want(+)\n%s\n", diff)
				}
				if diff := cmp.Diff(got.Bytes(), want.Bytes()); diff != "" {
					t.Errorf("resutl: got(-) want(+)\n%s\n", diff)
				}
			}

			t.Run("format json", func(t *testing.T) {
				testFunc(t, true)
			})
			t.Run("raw json", func(t *testing.T) {
				testFunc(t, false)
			})
		})
	}
}
