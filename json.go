package main

import (
	"encoding/json"
	"fmt"
	"io"
)

type jsonOption struct {
	path   string
	stType string
	prefix string
	indent string
}

func newJsonOption(path, stType string, noformat bool) jsonOption {
	op := jsonOption{
		path:   path,
		stType: stType,
	}
	if !noformat {
		op.indent = "  "
	}
	return op
}

func stType2Json(writer io.Writer, op jsonOption) error {
	defineMap, err := stType2Map(op.path, op.stType)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(writer)
	encoder.SetIndent(op.prefix, op.indent)
	if err = encoder.Encode(defineMap); err != nil {
		return fmt.Errorf("write json: %w", err)
	}
	return nil
}
