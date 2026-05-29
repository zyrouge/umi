package utils

import (
	"encoding/json"
	"strings"
)

func EncodeJsonToReader(data any) (*strings.Reader, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return strings.NewReader(string(bytes)), nil
}
