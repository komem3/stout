package main

import "encoding/json"

type basicValue struct {
	v interface{}
}

func (b basicValue) basic() {}

var _ json.Marshaler = basicValue{}

func (b basicValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.v)
}

func isBasic(i interface{}) (basicValue, bool) {
	v, ok := i.(basicValue)
	return v, ok
}

var basicMap = map[string]basicValue{
	"string":  {"string"},
	"int":     {1},
	"int8":    {1},
	"int16":   {1},
	"int32":   {1},
	"int64":   {1},
	"uint":    {1},
	"uint8":   {1},
	"uint16":  {1},
	"uint32":  {1},
	"uint64":  {1},
	"float32": {1.1},
	"float64": {1.1},
	"bool":    {true},
	"byte":    {1},
	"rune":    {1},
	"uintptr": {1},
	"error":   {nil},
}
