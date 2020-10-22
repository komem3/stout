package main

import (
	"bytes"
	"encoding/json"
	"regexp"
)

type structDefine struct {
	field  string
	typ    string
	tag    string
	define interface{}
}

type orderDefine []structDefine

var _ json.Marshaler = (orderDefine)(nil)

var jsonTag = regexp.MustCompile(`json:"([^,"]+)(,[^,"]*)?(,[^,"]*)?"`)

var stringTag = map[string]string{
	"string":  "string",
	"int":     "1",
	"int8":    "1",
	"int16":   "1",
	"int32":   "1",
	"int64":   "1",
	"float32": "1.0",
	"float64": "1.0",
	"bool":    "true",
}

func (o orderDefine) MarshalJSON() ([]byte, error) {
	w := new(bytes.Buffer)
	jw := json.NewEncoder(w)
	_, err := w.WriteRune('{')
	if err != nil {
		return nil, err
	}
	for i, d := range o {
		field := d.field
		define := d.define
		if d.tag != "" {
			matches := jsonTag.FindStringSubmatch(d.tag)
			if len(matches) == 4 {
				switch {
				case matches[1] == "-" && len(matches[2]) == 0:
					continue
				case matches[2] == ",string" || matches[3] == ",string":
					define = stringTag[d.typ]
					fallthrough
				default:
					field = matches[1]
				}
			}
		}
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
