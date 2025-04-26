package response

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Message string `json:"message"`
}

// JSON sends a JSON response with the provided status code and payload
func JSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Error sends a JSON error response with the provided status code and error message
func Error(w http.ResponseWriter, code int, message string) {
	JSON(w, code, ErrorResponse{Message: message})
}