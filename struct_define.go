package main

import (
	"bytes"
	"encoding/json"
)

type structDefine struct {
	field  string
	define interface{}
}

type orderDefine []structDefine

var _ json.Marshaler = (*orderDefine)(nil)

func (o orderDefine) MarshalJSON() ([]byte, error) {
	w := new(bytes.Buffer)
	jw := json.NewEncoder(w)
	_, err := w.WriteRune('{')
	if err != nil {
		return nil, err
	}
	for i, d := range o {
		_, err := w.WriteString("\"" + d.field + "\":")
		if err != nil {
			return nil, err
		}
		err = jw.Encode(d.define)
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
