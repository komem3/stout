package main

import (
	"encoding/json"
	"fmt"
	"io"
)

// TODO: json parse test
func stType2Json(path, stType string, writer io.Writer) error {
	defineMap, err := stType2Map(path, stType)
	if err != nil {
		return err
	}
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", "  ")
	if err = encoder.Encode(defineMap); err != nil {
		return fmt.Errorf("write json: %w", err)
	}
	return nil
}
