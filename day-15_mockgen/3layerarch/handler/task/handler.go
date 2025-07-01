package taskhandler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"3layerarch/models"
)

type Handler struct {
	Service TaskService
}

func New(service TaskService) *Handler {
	return &Handler{Service: service}
}

// CreateTask godoc
// @Summary Create a new task
// @Description Creates a task with the given JSON body
// @Tags tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task to create"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Router /task [post]
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

// GetTask godoc
// @Summary Get task by ID
// @Description Returns a task given its ID
// @Tags tasks
// @Produce json
// @Param id path int true "Task ID"
// @Success 200 {object} models.Task
// @Failure 400 {string} string "Bad Request"
// @Router /task/{id} [get]
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

// ViewTasks godoc
// @Summary List all tasks
// @Description Returns a list of all tasks
// @Tags tasks
// @Produce json
// @Success 200 {array} models.Task
// @Failure 500 {string} string "Internal Server Error"
// @Router /task [get]
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

// UpdateTask godoc
// @Summary Update task status
// @Description Marks a task as completed
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Router /task/{id} [put]
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

// DeleteTask godoc
// @Summary Delete task
// @Description Deletes a task by ID
// @Tags tasks
// @Param id path int true "Task ID"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Router /task/{id} [delete]
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
