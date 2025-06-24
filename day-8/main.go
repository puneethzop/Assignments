package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

// Task represents a to-do item.
type Task struct {
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

// TaskManager handles task storage and operations.
type TaskManager struct {
	tasks []*Task
}

// AddTask handles POST /task.
func (tm *TaskManager) addTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	var input Task
	err = json.Unmarshal(body, &input)

	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	tm.tasks = append(tm.tasks, &Task{
		Task:      input.Task,
		Completed: false,
	})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	message := fmt.Sprintf("Task '%s' added successfully", input.Task)
	resp, err := json.Marshal(message)

	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	if _, writeErr := w.Write(resp); writeErr != nil {
		fmt.Println("Failed to write response:", writeErr)
	}
}

// GetByID handles GET /task/{id}.
func (tm *TaskManager) getByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 0 || id >= len(tm.tasks) {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(tm.tasks[id])
	if err != nil {
		http.Error(w, "Failed to encode task", http.StatusInternalServerError)
		return
	}

	if _, writeErr := w.Write(data); writeErr != nil {
		fmt.Println("Failed to write response:", writeErr)
	}
}

// ViewAll handles GET /task.
func (tm *TaskManager) viewAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(tm.tasks)
	if err != nil {
		http.Error(w, "Failed to encode tasks", http.StatusInternalServerError)
		return
	}

	if _, writeErr := w.Write(data); writeErr != nil {
		fmt.Println("Failed to write response:", writeErr)
	}
}

// CompleteTask handles PATCH /task/{id}.
func (tm *TaskManager) completeTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 0 || id >= len(tm.tasks) {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	tm.tasks[id].Completed = true
	message := fmt.Sprintf("Task '%s' marked as completed", tm.tasks[id].Task)

	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	if _, writeErr := w.Write(resp); writeErr != nil {
		fmt.Println("Failed to write response:", writeErr)
	}
}

// DeleteTask handles DELETE /task/{id}.
func (tm *TaskManager) deleteTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 0 || id >= len(tm.tasks) {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	message := fmt.Sprintf("Task '%s' deleted successfully", tm.tasks[id].Task)
	tm.tasks = append(tm.tasks[:id], tm.tasks[id+1:]...)

	w.Header().Set("Content-Type", "application/json")

	resp, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

	if _, writeErr := w.Write(resp); writeErr != nil {
		fmt.Println("Failed to write response:", writeErr)
	}
}

func main() {
	tm := &TaskManager{}

	http.HandleFunc("POST /task", tm.addTask)
	http.HandleFunc("GET /task/{id}", tm.getByID)
	http.HandleFunc("GET /task", tm.viewAll)
	http.HandleFunc("PATCH /task/{id}", tm.completeTask)
	http.HandleFunc("DELETE /task/{id}", tm.deleteTask)

	server := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	fmt.Println("Server listening on http://localhost:8080")

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
