package respond

import (
	"encoding/json"
	"net/http"
)

type ErrorBody struct {
	Code    string `json:"code"`
	Details any    `json:"details,omitempty"`
}

type SuccessBody struct {
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func Problem(w http.ResponseWriter, status int, code, msg string, details any) {
	if status == 0 {
		status = http.StatusInternalServerError
	}
	writeJSON(w, status, map[string]any{
		"message": msg,
		"error":   ErrorBody{Code: code, Details: details},
	})
}

func OK(w http.ResponseWriter, status int, msg string, data any) {
	if status == 0 {
		status = http.StatusOK
	}
	writeJSON(w, status, SuccessBody{Message: msg, Data: data})
}

func NoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
