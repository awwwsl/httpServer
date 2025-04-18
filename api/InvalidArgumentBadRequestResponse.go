package api

import (
	"encoding/json"
	"httpServer/validation"
)

type InvalidArgumentBadRequestResponse struct {
	Errors map[string][]*validation.ValidateError `json:"errors,omitempty"`
}

func (i *InvalidArgumentBadRequestResponse) ToJson() ([]byte, error) {
	return json.Marshal(i)
}
