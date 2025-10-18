package responses

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func RespondJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

func RespondError(w http.ResponseWriter, statusCode int, message string) {
	RespondJSON(w, statusCode, ErrorResponse{Error: message})
}

func RespondSuccess(w http.ResponseWriter, statusCode int, data interface{}) {
	RespondJSON(w, statusCode, SuccessResponse{Data: data})
}
