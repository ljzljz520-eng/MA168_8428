package api

import (
	"encoding/json"
	"net/http"
)

type Envelope struct {
	Data  any               `json:"data,omitempty"`
	Error string            `json:"error,omitempty"`
	Meta  map[string]string `json:"meta,omitempty"`
}

func success(value any) Envelope {
	return Envelope{Data: value, Meta: map[string]string{"service": "bookstore-recommendation"}}
}

func encodeEnvelope(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(success(value))
}
