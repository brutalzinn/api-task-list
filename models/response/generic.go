package response_entities

import (
	"encoding/json"
	"net/http"
)

type GenericResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func GenericMessageError(w http.ResponseWriter, r *http.Request, message string) {
	resp := GenericResponse{
		Error:   true,
		Message: message,
		Data:    make([]any, 0),
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func GenericError(w http.ResponseWriter, r *http.Request) {
	resp := GenericResponse{
		Error:   true,
		Message: "Generic error",
		Data:    make([]any, 0),
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func GenericOK(w http.ResponseWriter, r *http.Request, data any) {
	resp := GenericResponse{
		Error:   false,
		Message: "",
		Data:    data,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func GenericMessageOK(w http.ResponseWriter, r *http.Request, data any, message string) {
	resp := GenericResponse{
		Error:   false,
		Message: message,
		Data:    data,
	}
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
