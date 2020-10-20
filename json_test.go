package main

import (
	"bytes"
	"encoding/json"
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
			outJson   []byte
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
				path:   "./struct_test.go",
				stType: "SampleJson",
			},
			want{
				outStruct: SampleJson{
					ID:      1,
					Content: "string",
					PStr:    pStr("string"),
					Pointer: &Embbading{
						ID: "string",
					},
					Struct: Embbading{
						ID: "string",
					},
					Array: []string{"string"},
					ArraySct: []Embbading{
						{
							ID: "string",
						},
					},
					ArrayPSct: []*Embbading{
						{
							ID: "string",
						},
					},
					Uinter: 1,
					Embbading: Embbading{
						ID: "string",
					},
					PEmb: &PEmb{
						Content: "string",
					},
					private: "",
				},
				outJson: []byte(""),
				err:     nil,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var (
				got  = new(bytes.Buffer)
				want = new(bytes.Buffer)
			)

			encoder := json.NewEncoder(want)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(tt.want.outStruct); err != nil {
				t.Fatal(err)
			}

			err := stType2Json(tt.given.path, tt.given.stType, got)
			if diff := cmp.Diff(err, tt.want.err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("error: got(-) want(+)\n%s\n", diff)
			}
			if diff := cmp.Diff(got.Bytes(), want.Bytes()); diff != "" {
				t.Errorf("resutl: got(-) want(+)\n%s\n", diff)
			}
		})
	}
}
