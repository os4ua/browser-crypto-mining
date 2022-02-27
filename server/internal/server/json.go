package server

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Error errorDetails `json:"error"`
}

type errorDetails struct {
	Message string `json:"message"`
}

func respondJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(status)

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	err := enc.Encode(v)
	if err != nil {
		log.Panic(err)
	}
}

func newErrorResponse(msg string) errorResponse {
	return errorResponse{
		Error: errorDetails{
			Message: msg,
		},
	}
}
