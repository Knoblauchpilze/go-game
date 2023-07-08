package db

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Convertible interface {
	Convert() interface{}
}

func argToStr(arg interface{}) (string, error) {
	var raw []byte
	var out string
	var err error

	if convertible, ok := arg.(Convertible); ok {
		raw, err = json.Marshal(convertible.Convert())
		if err == nil {
			out = string(raw)
		}
	} else if str, ok := arg.(string); ok {
		out = str
	} else if id, ok := arg.(uuid.UUID); ok {
		out = id.String()
	} else {
		raw, err = json.Marshal(arg)
		if err == nil {
			out = string(raw)
		}
	}

	return out, err
}
