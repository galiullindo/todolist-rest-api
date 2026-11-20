package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *handler) respondJSON(w http.ResponseWriter, v any, code int) {
	headers := w.Header()
	headers.Set("Content-Type", "application/json; charset=utf-8")
	headers.Set("X-Content-Type-Options", "nosniff")

	if v == nil {
		w.WriteHeader(code)
		return
	}

	body, err := json.Marshal(v)
	if err != nil {
		log.Printf("respond json: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"internal server error"}`))
		return
	}

	w.WriteHeader(code)
	if _, err := w.Write(body); err != nil {
		log.Printf("respond json: %v\n", err)
	}
}

type ResponseError struct {
	Error string `json:"error"`
}

func (h *handler) respondJSONError(w http.ResponseWriter, err string, code int) {
	headers := w.Header()
	headers.Del("Content-Length")

	h.respondJSON(w, &ResponseError{Error: err}, code)
}
