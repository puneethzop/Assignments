package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	taskhandler "3layerarch/handler/task"
	userhandler "3layerarch/handler/user"

	taskservice "3layerarch/service/task"
	userservice "3layerarch/service/user"

	taskstore "3layerarch/store/task"
	userstore "3layerarch/store/user"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:root123@tcp(localhost:3306)/test_db")
	if err != nil {
		log.Fatal("DB connection error:", err)
	}
	defer db.Close()

	// Ensure TASKS table exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS TASKS (
		id INT AUTO_INCREMENT PRIMARY KEY,
		task TEXT,
		completed BOOL DEFAULT FALSE
	);`)
	if err != nil {
		log.Fatal("Failed to create TASKS table:", err)
	}

	// Ensure USERS table exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS USERS (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(100)
	);`)
	if err != nil {
		log.Fatal("Failed to create USERS table:", err)
	}

	// Task dependency setup
	taskStore := taskstore.New(db)
	taskService := taskservice.New(taskStore)
	taskHandler := taskhandler.New(taskService)

	// User dependency setup
	userStore := userstore.New(db)
	userService := userservice.New(userStore)
	userHandler := userhandler.New(userService)

	// Task routes
	http.HandleFunc("POST /task", taskHandler.CreateTask)
	http.HandleFunc("GET /task", taskHandler.ViewTasks)
	http.HandleFunc("GET /task/{id}", taskHandler.GetTask)
	http.HandleFunc("PUT /task/{id}", taskHandler.UpdateTask)
	http.HandleFunc("DELETE /task/{id}", taskHandler.DeleteTask)

	// User routes
	http.HandleFunc("POST /user", userHandler.CreateUser)
	http.HandleFunc("GET /user/{id}", userHandler.GetUser)

	// Server configuration
	srv := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	//log.Println("Server running on http://localhost:8080")
	log.Fatal(srv.ListenAndServe())
}
