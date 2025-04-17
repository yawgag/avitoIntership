package errorsHandl

import (
	"encoding/json"
	"net/http"
)

type jsonError struct {
	Message string `json:"message"`
}

func SendJsonError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(jsonError{
		Message: message,
	})
}
