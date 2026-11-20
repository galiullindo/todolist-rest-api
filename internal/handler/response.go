package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseError struct {
	Error string `json:"error"`
}

func (h *handler) RespondErrorJSON(w http.ResponseWriter, err string, code int) {
	headers := w.Header()
	headers.Del("Content-Length")
	headers.Set("Content-Type", "application/json; charset=utf-8")
	headers.Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(code)
	body := ResponseError{Error: err}
	if err := json.NewEncoder(w).Encode(&body); err != nil {
		log.Printf("respond error json: %v\n", err)
	}
}

func (h *handler) RespondJSON(w http.ResponseWriter, body any, code int) {
	headers := w.Header()
	headers.Set("Content-Type", "application/json; charset=utf-8")
	headers.Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(code)
	if body == nil {
		return
	}
	if err := json.NewEncoder(w).Encode(&body); err != nil {
		log.Printf("respond json: %v\n", err)
	}
}
