package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Task struct {
	Task      string
	Completed bool
}

var tasks []*Task

func addTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)

	tasks = append(tasks, &Task{string(body), false})
	w.WriteHeader(http.StatusCreated)
}

func getByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
		return
	}
	w.Write([]byte(tasks[id].Task))
	w.WriteHeader(http.StatusOK)
}

func viewTask(w http.ResponseWriter, r *http.Request) {
	for i, t := range tasks {
		//fmt.Fprintf(w, "ID: %d, Task: %s, Completed: %v\n", i, t.Task, t.Completed)
		_, _ = w.Write([]byte(fmt.Sprintf("ID: %d, Task: %s, Completed: %v\n", i, t.Task, t.Completed)))
	}
}

func completeTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
		return
	}
	tasks[id].Completed = true
	w.WriteHeader(http.StatusOK)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
		return
	}
	tasks = append(tasks[:id], tasks[id+1:]...)
	w.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("POST /task", addTask)
	http.HandleFunc("GET /task/{id}", getByID)
	http.HandleFunc("GET /task", viewTask)
	http.HandleFunc("PUT /task/{id}", completeTask)
	http.HandleFunc("DELETE /task/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Not able to start server")
	}
}
