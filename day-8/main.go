package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type Task struct {
	Task      string `json:"task"`
	Completed bool   `json:"completed"`
}

var tasks []*Task

func addTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	task := string(body)
	tasks = append(tasks, &Task{Task: task, Completed: false})

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	message := fmt.Sprintf("Task '%s' added successfully", task)
	resp, _ := json.Marshal(message)
	w.Write(resp)

}

func getByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 || id > len(tasks) {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	data, _ := json.Marshal(tasks[id])
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func viewTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	data, _ := json.Marshal(tasks)
	w.Write(data)
}

func completeTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 || id > len(tasks) {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	tasks[id].Completed = true

	message := fmt.Sprintf("Task '%s' marked as completed", tasks[id].Task)

	w.Header().Set("Content-Type", "application/json")
	resp, _ := json.Marshal(message)
	w.Write(resp)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id < 1 || id > len(tasks) {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}
	message := fmt.Sprintf("Task '%s' deleted successfully", tasks[id].Task)

	tasks = append(tasks[:id], tasks[id+1:]...)

	w.Header().Set("Content-Type", "application/json")
	resp, _ := json.Marshal(message)
	w.Write(resp)
}

func main() {

	http.HandleFunc("POST /task", addTask)
	http.HandleFunc("GET /task/{id}", getByID)
	http.HandleFunc("GET /task", viewTask)
	http.HandleFunc("PATCH /task/{id}", completeTask)
	http.HandleFunc("DELETE /task/{id}", deleteTask)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Not able to start server")
	}
}
