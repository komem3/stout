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
			out interface{}
			err error
		}
	)
	tests := []struct {
		name  string
		given given
		want  want
	}{
		{
			"built-in types",
			given{
				path:   "./define_test.go",
				stType: "BuiltInTypes",
			},
			want{
				out: BuiltInTypes{
					Bool:    true,
					Byte:    1,
					Error:   nil,
					Float32: 1.1,
					Float64: 1.1,
					Int16:   1,
					Int32:   1,
					Int64:   1,
					Int8:    1,
					Rune:    1,
					String:  "string",
					Uint:    1,
					Uint16:  1,
					Uint32:  1,
					Uint64:  1,
					Uint8:   1,
					Uintptr: 1,
				},
			},
		},
		{
			"Not Struct(int)",
			given{
				"./define_test.go",
				"Integer",
			},
			want{
				out: Integer(1),
			},
		},
		{
			"Not Struct(array)",
			given{
				"./define_test.go",
				"ArraySamePkg",
			},
			want{
				out: ArraySamePkg{
					{
						IDString: "string",
					},
				},
			},
		},
		{
			"SimpleStruct",
			given{
				path:   "./define_test.go",
				stType: "SimpleStruct",
			},
			want{
				out: SimpleStruct{
					pfield: 1,
					IntF:   1,
					PstrF:  func(str string) *PString { s := PString(str); return &s }("string"),
					SameF: SamePkg{
						IDString: "string",
					},
					PSameF: &PtrSamePkg{
						Content: "string",
						Next:    "string",
					},
					ArraySameF: []SamePkg{
						{
							IDString: "string",
						},
					},
					ArrayPtrSameF: []*PtrSamePkg{
						{
							Content: "string",
							Next:    "string",
						},
					},
					PtrArraySameF: &PtrArraySamePkg{
						{
							IDString: "string",
						},
					},
					OtherIntF:    1,
					PrtOtherStrF: func(str string) *testutil.PtrOtherString { s := testutil.PtrOtherString(str); return &s }("string"),
					OtherF: testutil.OtherPkg{
						PkgContent: "string",
					},
					PtrOtherF: &testutil.PtrOtherPkg{
						PkgContentDiff: "string",
					},
					OtherArryF: []testutil.OtherPkg{
						{
							PkgContent: "string",
						},
					},
					OtherArrayPtrF: []*testutil.PtrOtherPkg{
						{
							PkgContentDiff: "string",
						},
					},
					PrtOtherArrayF: &testutil.PtrOtherArraySamePkg{
						{
							PkgContent: "string",
						},
					},
					CombinedF: CombineCase{
						CombindContent: "string",
					},
					AstF: ast.ArrayType{
						Lbrack: 1,
						Len:    nil,
						Elt:    nil,
					},
				},
			},
		},
		{
			"Emmbaddings",
			given{
				path:   "./define_test.go",
				stType: "Emmbaddings",
			},
			want{
				out: Emmbaddings{
					int:     1,
					private: 1,
					Integer: 1,
					PString: func(str string) *PString { s := PString(str); return &s }("string"),
					SamePkg: SamePkg{
						IDString: "string",
					},
					PtrSamePkg: &PtrSamePkg{
						Content: "string",
						Next:    "string",
					},
					ArraySamePkg: []SamePkg{
						{
							IDString: "string",
						},
					},
					ArrayPtrSampePkg: []*PtrSamePkg{
						{
							Content: "string",
							Next:    "string",
						},
					},
					PtrArraySamePkg: &PtrArraySamePkg{
						{
							IDString: "string",
						},
					},
					OtherInteger:   1,
					PtrOtherString: func(str string) *testutil.PtrOtherString { s := testutil.PtrOtherString(str); return &s }("string"),
					OtherPkg: testutil.OtherPkg{
						PkgContent: "string",
					},
					PtrOtherPkg: &testutil.PtrOtherPkg{
						PkgContentDiff: "string",
					},
					OtherArraySamePkg: []testutil.OtherPkg{
						{
							PkgContent: "string",
						},
					},
					OtherArrayPtrSampePkg: []*testutil.PtrOtherPkg{
						{
							PkgContentDiff: "string",
						},
					},
					PtrOtherArraySamePkg: &testutil.PtrOtherArraySamePkg{
						{
							PkgContent: "string",
						},
					},
					CombineCase: CombineCase{
						CombindContent: "string",
					},
					ArrayType: ast.ArrayType{
						Lbrack: 1,
						Len:    nil,
						Elt:    nil,
					},
				},
			},
		},
		{
			"json tag",
			given{
				path:   "./define_test.go",
				stType: "JsonTag",
			},
			want{
				out: JsonTag{
					ID:      1,
					Content: "string",
					Tag:     1,
					Ignore:  1,
					Hyphen:  1,
					Bool:    true,
					Byte:    1,
					Error:   nil,
					Float32: 1.1,
					Float64: 1.1,
					Int16:   1,
					Int32:   1,
					Int64:   1,
					Int8:    1,
					Rune:    1,
					String:  "string",
					Uint:    1,
					Uint16:  1,
					Uint32:  1,
					Uint64:  1,
					Uint8:   1,
					Uintptr: 1,
				},
			},
		},
		{
			"duplicate fileds",
			given{
				path:   "./define_test.go",
				stType: "DuplicateFields",
			},
			want{
				out: DuplicateFields{
					Field:     1,
					Duplicate: "string",
				},
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

				if err := encoder.Encode(tt.want.out); err != nil {
					t.Fatal(err)
				}

				op := newJsonOption(tt.given.path, tt.given.stType, !format)
				err := stType2Json(got, op)
				if diff := cmp.Diff(err, tt.want.err, cmpopts.EquateErrors()); diff != "" {
					t.Errorf("error: got(-) want(+)\n%s\n", diff)
				}
				if diff := cmp.Diff(got.String(), want.String()); diff != "" {
					t.Log(got.String())
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
