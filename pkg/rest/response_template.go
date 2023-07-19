package rest

import (
	"encoding/json"

	"github.com/google/uuid"
)

type ResponseTemplate struct {
	RequestId uuid.UUID
	Status    string
	Details   json.RawMessage `json:",omitempty"`
}
