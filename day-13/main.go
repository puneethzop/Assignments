package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type input struct {
	data *sql.DB
}

func hellohandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write([]byte("Hello, World!")); err != nil {
		log.Printf("Error writing response: %v", err)
	}
}

func (db *input) addTask(w http.ResponseWriter, r *http.Request) {
	if err := r.Body.Close(); err != nil {
		log.Printf("Error closing request body: %v", err)
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading body", http.StatusBadRequest)
		return
	}

	var inputData struct {
		Task string `json:"task"`
	}

	if err := json.Unmarshal(body, &inputData); err != nil || inputData.Task == "" {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	_, err = db.data.Exec("INSERT INTO TASKS (task, completed) VALUES (?, ?)", inputData.Task, false)
	if err != nil {
		http.Error(w, "Failed to insert task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

}

func (db *input) getByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var i int
	var task string
	var completed bool

	err = db.data.QueryRow("SELECT id, task, completed FROM TASKS WHERE id = ?", id).Scan(&i, &task, &completed)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Task not found", http.StatusNotFound)
			return
		}
		http.Error(w, "DB error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "ID: %d, Task: %s, Completed: %v", id, task, completed)
}

func (db *input) viewTask(w http.ResponseWriter, _ *http.Request) {
	rows, err := db.data.Query("SELECT id, task, completed FROM TASKS")
	if err != nil {
		http.Error(w, "Failed to query tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	w.Header().Set("Content-Type", "text/plain")

	for rows.Next() {
		var id int
		var task string
		var completed bool

		if err := rows.Scan(&id, &task, &completed); err != nil {
			http.Error(w, "Failed to read task row", http.StatusInternalServerError)
			log.Printf("%v", err)
			return
		}

		if _, err := fmt.Fprintf(w, "ID: %d, Task: %s, Completed: %t\n", id, task, completed); err != nil {
			log.Printf("Error writing task row: %v", err)
			return
		}
	}
}

func (db *input) completeTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	res, err := db.data.Exec("UPDATE TASKS SET completed = true WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	affected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking update result", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (db *input) deleteTask(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	res, err := db.data.Exec("DELETE FROM TASKS WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	affected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking delete result", http.StatusInternalServerError)
		return
	}

	if affected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func main() {
	db := &input{}
	var err error

	db.data, err = sql.Open("mysql", "root:root123@tcp(localhost:3306)/test_db")
	if err != nil {
		log.Fatal("Error connecting to DB -> ", err)
	}
	defer db.data.Close()

	if err := db.data.Ping(); err != nil {
		log.Fatal("Error pinging DB -> ", err)
	}

	_, err = db.data.Exec("CREATE TABLE IF NOT EXISTS TASKS (id int auto_increment primary key, task text, completed bool);")
	if err != nil {
		log.Fatal("Error creating table -> ", err)
	}

	http.HandleFunc("/", hellohandler)
	http.HandleFunc("POST /task", db.addTask)
	http.HandleFunc("GET /task/{id}", db.getByID)
	http.HandleFunc("GET /task", db.viewTask)
	http.HandleFunc("PUT /task/{id}", db.completeTask)
	http.HandleFunc("DELETE /task/{id}", db.deleteTask)

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}
