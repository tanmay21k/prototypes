// Package json contains small helpers for JSON HTTP responses.
package json

import (
	"encoding/json"
	"net/http"
)

// Write encodes data as JSON and writes it to the provided response writer.
func Write(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
