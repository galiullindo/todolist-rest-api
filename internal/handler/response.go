package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type ResponseError struct {
	Error string `json:"error"`
}

func RespondErrorJSON(w http.ResponseWriter, err string, code int) {
	h := w.Header()
	h.Del("Content-Length")
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("X-Content-Type-Options", "nosniff")

	w.WriteHeader(code)
	body := ResponseError{Error: err}
	if err := json.NewEncoder(w).Encode(&body); err != nil {
		log.Printf("respond error json: %v\n", err)
	}
}
