// Package json for handling Json, encoding, decoding and also sending json response to client
package json

import (
	"encoding/json"
	"net/http"

	"github.com/odundlaw/cbt-backend/internal/jwt"
)

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type APIResponse struct {
	Success bool         `json:"success"`
	Message string       `json:"message"`
	Data    any          `json:"data,omitempty"`
	Token   *jwt.Token   `json:"token,omitempty"`
	Errors  []FieldError `json:"errors,omitempty"`
}

func JSONSuccess(w http.ResponseWriter, status int, message string, data any, token *jwt.Token) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
		Token:   token,
	}

	json.NewEncoder(w).Encode(res)
}

func JSONError(w http.ResponseWriter, status int, message string, errs []FieldError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	res := APIResponse{
		Success: false,
		Message: message,
		Errors:  errs,
	}
	json.NewEncoder(w).Encode(res)
}

func ReadJSON(r *http.Request, data any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}
