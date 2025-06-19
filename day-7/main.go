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

func AddTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, _ := io.ReadAll(r.Body)
	tasks = append(tasks, &Task{Task: string(body), Completed: false})
}

func getByID(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
		return
	}
	fmt.Fprintf(w, "%v\n", tasks[id])
}

func viewTask(w http.ResponseWriter, r *http.Request) {
	for i, t := range tasks {
		fmt.Fprintf(w, "ID: %d, Task: %s, Completed: %v\n", i, t.Task, t.Completed)
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
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
		return
	}
	tasks = append(tasks[:id], tasks[id+1:]...)
}

func main() {
	http.HandleFunc("POST /task", AddTask)
	http.HandleFunc("GET /task/{id}", getByID)
	http.HandleFunc("GET /task", viewTask)
	http.HandleFunc("PATCH /task/{id}", completeTask)
	http.HandleFunc("DELETE /task/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Not able to start server")
	}
}
