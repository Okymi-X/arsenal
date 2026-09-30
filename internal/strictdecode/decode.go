// Package strictdecode decodes persisted configuration while rejecting fields
// the current schema does not understand.
package strictdecode

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
)

// JSON decodes one JSON value and rejects unknown fields or trailing values.
func JSON(data []byte, target any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		if err == nil {
			return fmt.Errorf("multiple JSON values")
		}
		return err
	}
	return nil
}

// TOML decodes one TOML document and rejects keys absent from the target schema.
func TOML(data []byte, target any) error {
	metadata, err := toml.NewDecoder(bytes.NewReader(data)).Decode(target)
	if err != nil {
		return err
	}
	undecoded := metadata.Undecoded()
	if len(undecoded) != 0 {
		return fmt.Errorf("unsupported TOML key %q", undecoded[0].String())
	}
	return nil
}
