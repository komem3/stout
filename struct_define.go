package main

import (
	"bytes"
	"encoding/json"
	"regexp"
)

type fieldDefine struct {
	field  string
	typ    string
	tag    string
	define json.Marshaler
	enable bool
}

func (f fieldDefine) empty() bool {
	return !f.enable
}

type orderDefine []fieldDefine

type arrayDefines []json.Marshaler

var _ json.Marshaler = (orderDefine)(nil)

var _ json.Marshaler = (arrayDefines)(nil)

var jsonTag = regexp.MustCompile(`json:"([^,"]*)(,[^,"]*)?(,[^,"]*)?"`)

var stringTag = map[string]basicValue{
	"string":  {"\"string\""},
	"byte":    {"1"},
	"rune":    {"1"},
	"int":     {"1"},
	"int8":    {"1"},
	"int16":   {"1"},
	"int32":   {"1"},
	"int64":   {"1"},
	"uint":    {"1"},
	"uint8":   {"1"},
	"uint16":  {"1"},
	"uint32":  {"1"},
	"uint64":  {"1"},
	"float32": {"1.1"},
	"float64": {"1.1"},
	"bool":    {"true"},
	"uintptr": {"1"},
}

func (ds arrayDefines) MarshalJSON() ([]byte, error) {
	w := new(bytes.Buffer)
	jw := json.NewEncoder(w)
	_, err := w.WriteRune('[')
	if err != nil {
		return nil, err
	}
	for i, d := range ds {
		if err = jw.Encode(d); err != nil {
			return nil, err
		}
		if i == len(ds)-1 {
			continue
		}
		_, err = w.WriteRune(',')
		if err != nil {
			return nil, err
		}
	}
	_, err = w.WriteRune(']')
	if err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func allocArray() arrayDefines {
	return make(arrayDefines, 1)
}

func (o orderDefine) MarshalJSON() ([]byte, error) {
	existsMap := make(map[string]int)
	for i, d := range o {
		if d.tag != "" {
			matches := jsonTag.FindStringSubmatch(d.tag)
			if len(matches) == 4 {
				switch {
				case matches[1] == "-" && len(matches[2]) == 0:
					o[i].enable = false
				case matches[2] == ",string" || matches[3] == ",string":
					o[i].define = stringTag[d.typ]
					fallthrough
				default:
					if matches[1] != "" {
						o[i].field = matches[1]
					}
				}
			}
		}
		field := o[i].field
		existIndex, ok := existsMap[field]
		if ok {
			o[existIndex].enable = false
		}
		existsMap[field] = i
	}
	return notDuplicateMarshall(o)
}

func notDuplicateMarshall(o orderDefine) ([]byte, error) {
	w := new(bytes.Buffer)
	jw := json.NewEncoder(w)
	_, err := w.WriteRune('{')
	if err != nil {
		return nil, err
	}
	if len(o) == 1 && o[0].field == "" {
		return json.Marshal(o[0].define)
	}
	for i, d := range o {
		if d.empty() {
			continue
		}
		field := d.field
		define := d.define
		_, err := w.WriteString("\"" + field + "\":")
		if err != nil {
			return nil, err
		}
		err = jw.Encode(define)
		if err != nil {
			return nil, err
		}
		if i == len(o)-1 {
			continue
		}
		_, err = w.WriteRune(',')
		if err != nil {
			return nil, err
		}
	}
	_, err = w.WriteRune('}')
	if err != nil {
		return nil, err
	}
	return w.Bytes(), err
}
