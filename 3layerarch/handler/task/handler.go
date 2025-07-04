package taskhandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"3layerarch/models"
)

type TaskService interface {
	CreateTask(t models.Task) error
	GetTask(id int) (models.Task, error)
	ViewTasks() ([]models.Task, error)
	UpdateTask(id int) error
	DeleteTask(id int) error
}

type Handler struct {
	Service TaskService
}

func New(service TaskService) *Handler {
	return &Handler{Service: service}
}

func (h *Handler) CreateTask(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil || len(body) == 0 {
		http.Error(w, "Empty or unreadable body", http.StatusBadRequest)
		return
	}
	var t models.Task
	if err := json.Unmarshal(body, &t); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}
	if err := h.Service.CreateTask(t); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	t, err := h.Service.GetTask(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	b, _ := json.Marshal(t)
	if _, err := w.Write(b); err != nil {
		fmt.Printf("failed to write response: %v\n", err)
	}
}

func (h *Handler) ViewTasks(w http.ResponseWriter, _ *http.Request) {
	tasks, err := h.Service.ViewTasks()
	if err != nil {
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	b, _ := json.Marshal(tasks)
	if _, err := w.Write(b); err != nil {
		fmt.Printf("failed to write response: %v\n", err)
	}
}

func (h *Handler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	if err := h.Service.UpdateTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	if idStr == "" {
		http.Error(w, "Missing ID", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}
	if err := h.Service.DeleteTask(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
