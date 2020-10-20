package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestStType2Map(t *testing.T) {
	type (
		given struct {
			path   string
			stType string
		}
		want struct {
			define structDefine
			err    error
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
				define: structDefine{
					"ID":      1,
					"Content": "string",
					"PStr":    "string",
					"Pointer": structDefine{
						"ID": "string",
					},
					"Struct": structDefine{
						"ID": "string",
					},
					"Array": []interface{}{"string"},
					"ArraySct": []interface{}{
						structDefine{"ID": "string"},
					},
					"ArrayPSct": []interface{}{
						structDefine{"ID": "string"},
					},
					"Uinter": 1,
				},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := stType2Map(tt.given.path, tt.given.stType)
			if diff := cmp.Diff(err, tt.want.err, cmpopts.EquateErrors()); diff != "" {
				t.Errorf("error: got(-) want(+)\n%s\n", diff)
			}
			if diff := cmp.Diff(got, tt.want.define); diff != "" {
				t.Errorf("resutl: got(-) want(+)\n%s\n", diff)
			}
		})
	}
}
