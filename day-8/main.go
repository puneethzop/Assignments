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

	task := string(body)
	tm.tasks = append(tm.tasks, &Task{Task: task, Completed: false})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	message := fmt.Sprintf("Task '%s' added successfully", task)
	resp, err := json.Marshal(message)

	if err == nil {
		if _, err := w.Write(resp); err != nil {
			fmt.Println("Failed to write response:", err)
		}
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
	if err == nil {
		if _, err := w.Write(data); err != nil {
			fmt.Println("Failed to write response:", err)
		}
	}
}

// ViewAll handles GET /task.
func (tm *TaskManager) viewAll(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(tm.tasks)
	if err == nil {
		if _, err := w.Write(data); err != nil {
			fmt.Println("Failed to write response:", err)
		}
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
	if err == nil {
		if _, err := w.Write(resp); err != nil {
			fmt.Println("Failed to write response:", err)
		}
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
	if err == nil {
		if _, err := w.Write(resp); err != nil {
			fmt.Println("Failed to write response:", err)
		}
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

	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Failed to start server:", err)
	}
}
