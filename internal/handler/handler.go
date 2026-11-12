package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	errs "github.com/galiullindo/todolist-rest-api/internal/errors"
	"github.com/galiullindo/todolist-rest-api/internal/model"
)

type Service interface {
	Get() ([]model.Task, error)
	Check() error
	Create(title string, status string) (model.Task, error)

	GetByID(id int64) (model.Task, error)
	CheckExistenceByID(id int64) error
	ReplaceByID(id int64, title string, status string) (model.Task, error)
	DeleteByID(id int64) error
}

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h *handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /tasks", h.Get)
	mux.HandleFunc("HEAD /tasks", h.Check)
	mux.HandleFunc("POST /tasks", h.Create)

	mux.HandleFunc("GET /tasks/{id}", h.GetByID)
	mux.HandleFunc("HEAD /tasks/{id}", h.CheckByID)
	mux.HandleFunc("PUT /tasks/{id}", h.ReplaceByID)
	mux.HandleFunc("DELETE /tasks/{id}", h.DeleteByID)
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.Get()
	if err != nil {
		log.Printf("task handler: get: %v\n", err)
		http.Error(w, "interna server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Printf("task handler: get: %v\n", err)
	}
}

func (h *handler) Check(w http.ResponseWriter, r *http.Request) {
	err := h.service.Check()
	if err != nil {
		log.Printf("task handler: check: %v\n", err)
		http.Error(w, "interna server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Status string
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("task handler: create: %v\n", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	task, err := h.service.Create(input.Title, input.Status)
	if err != nil {
		log.Printf("task handler: create: %v\n", err)
		http.Error(w, "interna server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("task handler: create: %v\n", err)
	}
}

func (h *handler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("task handler: get by id: %v\n", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetByID(id)
	if err != nil {
		log.Printf("task handler: get by id: %v\n", err)
		if errors.Is(err, errs.ErrTaskNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}

		http.Error(w, "interna server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("task handler: get by id: %v\n", err)
	}
}

func (h *handler) CheckByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("task handler: check by id: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.service.CheckExistenceByID(id); err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		log.Printf("task handler: check by id: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

func (h *handler) ReplaceByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("task handler: replace by id: %v\n", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var input struct {
		Title  string
		Status string
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("task handler: replace by id: %v\n", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	task, err := h.service.ReplaceByID(id, input.Title, input.Status)
	if err != nil {
		log.Printf("task handler: replace by id: %v\n", err)
		if errors.Is(err, errs.ErrTaskNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Printf("task handler: replace by id: %v\n", err)
	}
}

func (h *handler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("task handler: delete by id: %v\n", err)
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteByID(id); err != nil {
		log.Printf("task handler: replace by id: %v\n", err)
		if errors.Is(err, errs.ErrTaskNotFound) {
			http.Error(w, "task not found", http.StatusNotFound)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
